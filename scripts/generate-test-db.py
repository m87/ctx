#!/usr/bin/env python3
from __future__ import annotations

import argparse
import sqlite3
from datetime import datetime, time, timedelta, timezone
from pathlib import Path
from uuid import NAMESPACE_URL, uuid5
from zoneinfo import ZoneInfo


DATABASE_VERSION = "0.5.0"
TEST_ID_NAMESPACE = "ctx2-test-database-v1"
TIMEZONE_NAME = "Europe/Warsaw"
NANOSECONDS_PER_SECOND = 1_000_000_000
SEEDED_DAYS = 7


def stable_id(*parts: str) -> str:
    return str(uuid5(NAMESPACE_URL, "/".join((TEST_ID_NAMESPACE, *parts))))


def sql_time(value: datetime) -> str:
    utc = value.astimezone(timezone.utc)
    return utc.strftime("%Y-%m-%d %H:%M:%S.%f") + "000+00:00"


def duration_ns(duration: timedelta) -> int:
    return int(duration.total_seconds() * NANOSECONDS_PER_SECOND)


def local_zone() -> ZoneInfo:
    return ZoneInfo(TIMEZONE_NAME)


def recent_day_start(now: datetime, days_back: int, hour: int, minute: int = 0) -> datetime:
    target_date = now.astimezone(local_zone()).date() - timedelta(days=days_back)
    return datetime.combine(target_date, time(hour, minute), tzinfo=local_zone())


def create_schema(conn: sqlite3.Connection) -> None:
    conn.executescript(
        """
        CREATE TABLE node_cores (
          id char(36),
          namespace_id char(36),
          parent_id char(36),
          kind text NOT NULL DEFAULT "",
          status text NOT NULL DEFAULT "",
          name text NOT NULL,
          created_at datetime NOT NULL,
          updated_at datetime NOT NULL,
          PRIMARY KEY (id)
        );
        CREATE INDEX idx_node_cores_name ON node_cores(name);
        CREATE INDEX idx_node_cores_status ON node_cores(status);
        CREATE INDEX idx_node_cores_kind ON node_cores(kind);
        CREATE INDEX idx_node_cores_parent_id ON node_cores(parent_id);
        CREATE INDEX idx_parent_id ON node_cores(parent_id);
        CREATE INDEX idx_node_cores_namespace_id ON node_cores(namespace_id);
        CREATE INDEX idx_namespace_id ON node_cores(namespace_id);

        CREATE TABLE tags (
          id char(36),
          namespace_id char(36),
          name text NOT NULL,
          created_at datetime NOT NULL,
          PRIMARY KEY (id)
        );
        CREATE INDEX idx_tags_name ON tags(name);
        CREATE INDEX idx_tags_namespace_id ON tags(namespace_id, namespace_id);

        CREATE TABLE node_tags (
          node_id char(36),
          tag_id char(36),
          PRIMARY KEY (node_id, tag_id)
        );
        CREATE INDEX idx_node_tag ON node_tags(node_id, tag_id);

        CREATE TABLE kvs (
          node_id char(36),
          key text,
          value_text text,
          value_number real,
          value_int integer,
          value_int64 bigint,
          value_bool boolean,
          value_time datetime,
          PRIMARY KEY (node_id, key)
        );
        CREATE INDEX idx_kv_key ON kvs(key);
        CREATE INDEX idx_kv_node_id ON kvs(node_id);

        CREATE TABLE contents (
          node_id char(36),
          key text,
          value text,
          created_at datetime NOT NULL,
          updated_at datetime NOT NULL,
          PRIMARY KEY (node_id, key)
        );
        CREATE INDEX idx_content_key ON contents(key);
        CREATE INDEX idx_content_node_id ON contents(node_id);
        """
    )


def insert_node(
    conn: sqlite3.Connection,
    *,
    node_id: str,
    kind: str,
    name: str,
    namespace_id: str | None = None,
    parent_id: str | None = None,
    status: str = "",
    now: datetime,
) -> None:
    timestamp = sql_time(now)
    conn.execute(
        """
        INSERT INTO node_cores
          (id, namespace_id, parent_id, kind, status, name, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        """,
        (node_id, namespace_id, parent_id, kind, status, name, timestamp, timestamp),
    )


def insert_content(
    conn: sqlite3.Connection,
    *,
    node_id: str,
    key: str,
    value: str,
    now: datetime,
) -> None:
    timestamp = sql_time(now)
    conn.execute(
        """
        INSERT INTO contents (node_id, key, value, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
        """,
        (node_id, key, value, timestamp, timestamp),
    )


def insert_kv_text(conn: sqlite3.Connection, node_id: str, key: str, value: str) -> None:
    conn.execute(
        "INSERT INTO kvs (node_id, key, value_text) VALUES (?, ?, ?)",
        (node_id, key, value),
    )


def insert_kv_int64(conn: sqlite3.Connection, node_id: str, key: str, value: int) -> None:
    conn.execute(
        "INSERT INTO kvs (node_id, key, value_int64) VALUES (?, ?, ?)",
        (node_id, key, value),
    )


def insert_kv_time(conn: sqlite3.Connection, node_id: str, key: str, value: datetime) -> None:
    conn.execute(
        "INSERT INTO kvs (node_id, key, value_time) VALUES (?, ?, ?)",
        (node_id, key, sql_time(value)),
    )


def create_system_records(conn: sqlite3.Connection, now: datetime) -> None:
    insert_node(
        conn,
        node_id="systemInfoV1",
        kind="system_info",
        name="systemInfoV1",
        now=now,
    )
    insert_kv_text(conn, "systemInfoV1", "database_version", DATABASE_VERSION)

    insert_node(
        conn,
        node_id="settingsV1",
        kind="settings",
        name="settingsV1",
        now=now,
    )
    insert_kv_text(conn, "settingsV1", "client.general.firstDay", "Monday")
    insert_kv_text(conn, "settingsV1", "client.general.theme", "dark")


def create_workspace(
    conn: sqlite3.Connection,
    *,
    slug: str,
    name: str,
    description: str,
    now: datetime,
) -> str:
    workspace_id = stable_id("workspace", slug)
    insert_node(conn, node_id=workspace_id, kind="workspace", name=name, now=now)
    insert_content(
        conn,
        node_id=workspace_id,
        key="description",
        value=description,
        now=now,
    )
    return workspace_id


def create_context(
    conn: sqlite3.Connection,
    *,
    slug: str,
    name: str,
    workspace_id: str,
    description: str,
    now: datetime,
    status: str = "inactive",
) -> str:
    context_id = stable_id("context", slug)
    insert_node(
        conn,
        node_id=context_id,
        namespace_id=workspace_id,
        parent_id="",
        kind="context",
        status=status,
        name=name,
        now=now,
    )
    insert_content(
        conn,
        node_id=context_id,
        key="description",
        value=description,
        now=now,
    )
    return context_id


def create_interval(
    conn: sqlite3.Connection,
    *,
    slug: str,
    context_id: str,
    workspace_id: str,
    start: datetime,
    duration: timedelta,
    now: datetime,
    status: str = "completed",
) -> str:
    interval_id = stable_id("interval", slug)
    end = start + duration
    insert_node(
        conn,
        node_id=interval_id,
        namespace_id=workspace_id,
        parent_id=context_id,
        kind="interval",
        status=status,
        name=interval_id,
        now=now,
    )
    insert_kv_time(conn, interval_id, "start", start)
    insert_kv_text(conn, interval_id, "start_timezone", TIMEZONE_NAME)
    insert_kv_time(conn, interval_id, "end", end)
    insert_kv_text(conn, interval_id, "end_timezone", TIMEZONE_NAME)
    insert_kv_int64(conn, interval_id, "duration", duration_ns(duration))
    return interval_id


def create_interval_record(
    conn: sqlite3.Connection,
    *,
    slug: str,
    context_id: str,
    workspace_id: str,
    start: datetime | None,
    end: datetime | None,
    now: datetime,
    status: str,
) -> str:
    interval_id = stable_id("interval", slug)
    insert_node(
        conn,
        node_id=interval_id,
        namespace_id=workspace_id,
        parent_id=context_id,
        kind="interval",
        status=status,
        name=interval_id,
        now=now,
    )
    if start is not None:
        insert_kv_time(conn, interval_id, "start", start)
        insert_kv_text(conn, interval_id, "start_timezone", TIMEZONE_NAME)
    if end is not None:
        insert_kv_time(conn, interval_id, "end", end)
        insert_kv_text(conn, interval_id, "end_timezone", TIMEZONE_NAME)
    if start is not None and end is not None and end > start:
        insert_kv_int64(conn, interval_id, "duration", duration_ns(end - start))
    else:
        insert_kv_int64(conn, interval_id, "duration", 0)
    return interval_id


def seed_large_distribution_workspace(
    conn: sqlite3.Connection,
    *,
    micro_contexts: int,
    now: datetime,
) -> str:
    workspace_id = create_workspace(
        conn,
        slug="large-distribution",
        name="Large Distribution Workspace",
        description=(
            "Stress-test workspace with many contexts. Most micro contexts stay below "
            "one percent so the distribution chart can group them."
        ),
        now=now,
    )

    major_contexts = [
        ("deep-work", "Deep Work", 120, 8, 0),
        ("meetings", "Meetings", 75, 10, 15),
        ("product-planning", "Product Planning", 60, 11, 45),
        ("engineering-support", "Engineering Support", 45, 13, 0),
        ("research", "Research", 35, 14, 0),
        ("code-review", "Code Review", 30, 14, 45),
        ("administration", "Administration", 20, 15, 30),
        ("learning", "Learning", 15, 16, 0),
    ]

    for context_index, (slug, name, minutes, hour, minute) in enumerate(major_contexts):
        context_id = create_context(
            conn,
            slug=f"large-{slug}",
            name=name,
            workspace_id=workspace_id,
            description="Large workspace primary context.",
            now=now,
        )
        for days_back in range(SEEDED_DAYS):
            create_interval(
                conn,
                slug=f"large-{slug}-day-{days_back}",
                context_id=context_id,
                workspace_id=workspace_id,
                start=recent_day_start(now, days_back, hour, minute),
                duration=timedelta(minutes=minutes + ((context_index + days_back) % 3) * 5),
                now=now,
            )

    micro_names = [
        "Bug Triage",
        "Inbox Cleanup",
        "Release Note",
        "Design Ping",
        "Metrics Check",
        "Standup Follow-up",
        "Dependency Review",
        "Customer Note",
    ]
    for index in range(1, micro_contexts + 1):
        name = f"{micro_names[(index - 1) % len(micro_names)]} {index:02d}"
        context_id = create_context(
            conn,
            slug=f"large-micro-{index:02d}",
            name=name,
            workspace_id=workspace_id,
            description="Tiny context intentionally kept below one percent.",
            now=now,
        )
        create_interval(
            conn,
            slug=f"large-micro-{index:02d}-interval",
            context_id=context_id,
            workspace_id=workspace_id,
            start=recent_day_start(now, (index - 1) % SEEDED_DAYS, 17)
            + timedelta(minutes=((index - 1) // SEEDED_DAYS) * 6),
            duration=timedelta(minutes=2 + (index % 4)),
            now=now,
        )

    return workspace_id


def seed_small_healthy_workspace(conn: sqlite3.Connection, *, now: datetime) -> str:
    workspace_id = create_workspace(
        conn,
        slug="small-healthy",
        name="Small Healthy Workspace",
        description="Compact valid workspace with a few clean contexts and intervals.",
        now=now,
    )
    contexts = [
        ("writing", "Writing", [(0, 19, 0, 45), (2, 19, 0, 75), (5, 19, 0, 60)]),
        ("review", "Review", [(1, 19, 0, 30), (4, 19, 0, 45)]),
        ("administration", "Administration", [(0, 20, 0, 20), (6, 19, 45, 30)]),
    ]
    for slug, name, intervals in contexts:
        context_id = create_context(
            conn,
            slug=f"small-{slug}",
            name=name,
            workspace_id=workspace_id,
            description="Healthy sample context.",
            now=now,
        )
        for index, (days_back, hour, minute, minutes) in enumerate(intervals):
            create_interval(
                conn,
                slug=f"small-{slug}-interval-{index}",
                context_id=context_id,
                workspace_id=workspace_id,
                start=recent_day_start(now, days_back, hour, minute),
                duration=timedelta(minutes=minutes),
                now=now,
            )
    return workspace_id


def seed_integrity_error_workspace(
    conn: sqlite3.Connection,
    *,
    mismatch_workspace_id: str,
    now: datetime,
) -> str:
    workspace_id = create_workspace(
        conn,
        slug="integrity-errors",
        name="Integrity Error Workspace",
        description="Workspace with intentional broken records for the Data integrity view.",
        now=now,
    )
    start = recent_day_start(now, 1, 21)

    anchor_context_id = create_context(
        conn,
        slug="integrity-anchor",
        name="Integrity Anchor Context",
        workspace_id=workspace_id,
        description="Valid context used by intentionally broken intervals.",
        now=now,
    )
    create_interval(
        conn,
        slug="integrity-anchor-valid-interval",
        context_id=anchor_context_id,
        workspace_id=workspace_id,
        start=start,
        duration=timedelta(minutes=30),
        now=now,
    )

    create_context(
        conn,
        slug="broken-context-missing-workspace",
        name="Broken Context Missing Workspace",
        workspace_id="",
        description="Intentional issue: context has no workspace assigned.",
        now=now,
    )
    create_context(
        conn,
        slug="broken-context-missing-workspace-reference",
        name="Broken Context Missing Workspace Reference",
        workspace_id="missing-workspace-for-context",
        description="Intentional issue: context references a workspace that does not exist.",
        now=now,
    )

    create_interval(
        conn,
        slug="broken-interval-missing-context",
        context_id="",
        workspace_id=workspace_id,
        start=start + timedelta(hours=1),
        duration=timedelta(minutes=15),
        now=now,
    )
    create_interval(
        conn,
        slug="broken-interval-context-not-found",
        context_id="missing-context-for-interval",
        workspace_id=workspace_id,
        start=start + timedelta(hours=2),
        duration=timedelta(minutes=20),
        now=now,
    )
    create_interval(
        conn,
        slug="broken-interval-missing-workspace",
        context_id=anchor_context_id,
        workspace_id="",
        start=start + timedelta(hours=3),
        duration=timedelta(minutes=10),
        now=now,
    )
    create_interval(
        conn,
        slug="broken-interval-workspace-mismatch",
        context_id=anchor_context_id,
        workspace_id=mismatch_workspace_id,
        start=start + timedelta(hours=4),
        duration=timedelta(minutes=12),
        now=now,
    )
    create_interval(
        conn,
        slug="broken-interval-context-and-workspace-not-found",
        context_id="missing-context-and-workspace-for-interval",
        workspace_id="missing-workspace-for-interval",
        start=start + timedelta(hours=5),
        duration=timedelta(minutes=18),
        now=now,
    )

    create_interval_record(
        conn,
        slug="broken-completed-interval-missing-start",
        context_id=anchor_context_id,
        workspace_id=workspace_id,
        start=None,
        end=start + timedelta(hours=6, minutes=15),
        status="completed",
        now=now,
    )
    create_interval_record(
        conn,
        slug="broken-completed-interval-missing-end",
        context_id=anchor_context_id,
        workspace_id=workspace_id,
        start=start + timedelta(hours=6, minutes=30),
        end=None,
        status="completed",
        now=now,
    )

    active_context_a_id = create_context(
        conn,
        slug="broken-active-context-a",
        name="Broken Active Context A",
        workspace_id=workspace_id,
        description="Intentional issue: more than one context is active.",
        status="active",
        now=now,
    )
    create_interval(
        conn,
        slug="broken-active-context-a-ended-interval",
        context_id=active_context_a_id,
        workspace_id=workspace_id,
        start=start + timedelta(hours=7),
        duration=timedelta(minutes=10),
        now=now,
    )

    active_context_b_id = create_context(
        conn,
        slug="broken-active-context-b",
        name="Broken Active Context B",
        workspace_id=workspace_id,
        description="Intentional issue: more than one context is active.",
        status="active",
        now=now,
    )
    create_interval(
        conn,
        slug="broken-active-context-b-ended-interval",
        context_id=active_context_b_id,
        workspace_id=workspace_id,
        start=start + timedelta(hours=7, minutes=15),
        duration=timedelta(minutes=10),
        now=now,
    )

    active_interval_with_end_context_id = create_context(
        conn,
        slug="broken-active-interval-with-end-context",
        name="Broken Active Interval With End Context",
        workspace_id=workspace_id,
        description="Intentional issue: active interval has an end time.",
        status="active",
        now=now,
    )
    create_interval_record(
        conn,
        slug="broken-active-interval-with-end",
        context_id=active_interval_with_end_context_id,
        workspace_id=workspace_id,
        start=start + timedelta(hours=7, minutes=30),
        end=start + timedelta(hours=7, minutes=40),
        status="active",
        now=now,
    )

    return workspace_id


def validate_no_overlapping_intervals(conn: sqlite3.Connection) -> None:
    overlaps = conn.execute(
        """
        WITH intervals AS (
          SELECT
            core.id,
            core.name,
            start_kv.value_time AS start_time,
            end_kv.value_time AS end_time
          FROM node_cores core
          JOIN kvs start_kv ON start_kv.node_id = core.id AND start_kv.key = 'start'
          JOIN kvs end_kv ON end_kv.node_id = core.id AND end_kv.key = 'end'
          WHERE core.kind = 'interval'
        )
        SELECT
          first_interval.id,
          second_interval.id,
          first_interval.start_time,
          first_interval.end_time,
          second_interval.start_time,
          second_interval.end_time
        FROM intervals first_interval
        JOIN intervals second_interval ON first_interval.id < second_interval.id
        WHERE first_interval.start_time < second_interval.end_time
          AND second_interval.start_time < first_interval.end_time
        LIMIT 5
        """
    ).fetchall()

    if overlaps:
        formatted = "; ".join(
            f"{first_id} ({first_start}–{first_end}) overlaps "
            f"{second_id} ({second_start}–{second_end})"
            for first_id, second_id, first_start, first_end, second_start, second_end in overlaps
        )
        raise RuntimeError(f"Generated overlapping intervals: {formatted}")


def generate_database(
    output: Path,
    *,
    micro_contexts: int,
    include_integrity_errors: bool,
    force: bool,
) -> None:
    if output.exists() and not force:
        raise SystemExit(f"{output} already exists. Use --force to replace it.")

    output.parent.mkdir(parents=True, exist_ok=True)
    if output.exists():
        output.unlink()

    now = datetime.now(local_zone()).replace(microsecond=0)
    with sqlite3.connect(output) as conn:
        conn.execute("PRAGMA foreign_keys = OFF")
        create_schema(conn)
        create_system_records(conn, now)
        seed_large_distribution_workspace(conn, micro_contexts=micro_contexts, now=now)
        small_workspace_id = seed_small_healthy_workspace(conn, now=now)
        if include_integrity_errors:
            seed_integrity_error_workspace(
                conn,
                mismatch_workspace_id=small_workspace_id,
                now=now,
            )
        validate_no_overlapping_intervals(conn)
        conn.commit()


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Generate a ctx SQLite database for UI and data integrity testing.",
    )
    parser.add_argument(
        "-o",
        "--output",
        default="ctx.test.db",
        type=Path,
        help="Output database path. Defaults to ctx.test.db.",
    )
    parser.add_argument(
        "--micro-contexts",
        default=80,
        type=int,
        help="Number of sub-1%% contexts in Large Distribution Workspace.",
    )
    parser.add_argument(
        "--force",
        action="store_true",
        help="Replace the output database if it already exists.",
    )
    parser.add_argument(
        "--include-integrity-errors",
        action="store_true",
        help="Include an intentionally broken workspace for the Data integrity view.",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    if args.micro_contexts < 0:
        raise SystemExit("--micro-contexts cannot be negative.")

    generate_database(
        args.output,
        micro_contexts=args.micro_contexts,
        include_integrity_errors=args.include_integrity_errors,
        force=args.force,
    )
    print(f"Generated test database: {args.output}")
    print("Workspaces:")
    print("- Large Distribution Workspace")
    print("- Small Healthy Workspace")
    if args.include_integrity_errors:
        print("- Integrity Error Workspace")


if __name__ == "__main__":
    main()
