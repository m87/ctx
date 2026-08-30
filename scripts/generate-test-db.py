#!/usr/bin/env python3
"""Generate a ctx test database in the current relational 0.7.1 format."""

from __future__ import annotations

import argparse
import json
import sqlite3
from datetime import datetime, time, timedelta, timezone
from pathlib import Path
from uuid import NAMESPACE_URL, uuid4, uuid5
from zoneinfo import ZoneInfo


DATABASE_VERSION = "0.7.1"
TEST_ID_NAMESPACE = "ctx2-test-database-v1"
SEED_TIMEZONE_NAME = "UTC"
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
    return ZoneInfo(SEED_TIMEZONE_NAME)


def recent_day_start(now: datetime, days_back: int, hour: int, minute: int = 0) -> datetime:
    target_date = now.astimezone(local_zone()).date() - timedelta(days=days_back)
    return datetime.combine(target_date, time(hour, minute), tzinfo=local_zone())


def create_schema(conn: sqlite3.Connection) -> None:
    conn.executescript(
        """
        CREATE TABLE properties (
          id text PRIMARY KEY,
          database_version text,
          client_id text
        );

        CREATE TABLE client_properties (
          id text PRIMARY KEY,
          theme text,
          first_day text,
          timezone text,
          "values" text
        );

        CREATE TABLE workspaces (
          id text PRIMARY KEY,
          name text,
          description text
        );

        CREATE TABLE workspace_link_rules (
          workspace_id text,
          position integer,
          regexp text NOT NULL,
          link text NOT NULL,
          PRIMARY KEY (workspace_id, position),
          CONSTRAINT fk_workspaces_link_rules FOREIGN KEY (workspace_id)
            REFERENCES workspaces(id)
            ON DELETE CASCADE
        );

        CREATE TABLE projects (
          id text PRIMARY KEY,
          name text,
          parent_id text,
          workspace_id text NOT NULL
        );
        CREATE INDEX idx_projects_workspace_id ON projects(workspace_id);
        CREATE INDEX idx_projects_parent_id ON projects(parent_id);

        CREATE TABLE contexts (
          id text PRIMARY KEY,
          name text NOT NULL,
          workspace_id text NOT NULL,
          status text NOT NULL,
          archived numeric NOT NULL,
          description text,
          project_id text,
          CONSTRAINT fk_contexts_project_metadata FOREIGN KEY (project_id)
            REFERENCES projects(id)
        );
        CREATE INDEX idx_contexts_project_id ON contexts(project_id);
        CREATE INDEX idx_contexts_workspace_id ON contexts(workspace_id);

        CREATE TABLE tag (
          id text PRIMARY KEY,
          name text NOT NULL
        );

        CREATE TABLE context_tags (
          context_id text,
          tag_id text,
          PRIMARY KEY (context_id, tag_id),
          CONSTRAINT fk_context_tags_context_entity FOREIGN KEY (context_id)
            REFERENCES contexts(id),
          CONSTRAINT fk_context_tags_tag_entity FOREIGN KEY (tag_id)
            REFERENCES tag(id)
        );

        CREATE TABLE intervals (
          id text PRIMARY KEY,
          context_id text NOT NULL,
          start datetime NOT NULL,
          end datetime,
          duration integer NOT NULL,
          status text NOT NULL,
          workspace_id text NOT NULL
        );
        CREATE INDEX idx_intervals_workspace_id ON intervals(workspace_id);
        CREATE INDEX idx_intervals_status ON intervals(status);
        CREATE INDEX idx_intervals_context_id ON intervals(context_id);
        """
    )


def create_system_records(conn: sqlite3.Connection, now: datetime) -> None:
    client_id = str(uuid4())
    conn.execute(
        "INSERT INTO properties (id, database_version, client_id) VALUES (?, ?, ?)",
        ("system", DATABASE_VERSION, client_id),
    )
    client_values = {
        "client.general.theme": "dark",
        "client.general.firstDay": "Monday",
        "client.general.timeZone": "browser",
        "client.generator.sample": "preserved",
    }
    conn.execute(
        """
        INSERT INTO client_properties (id, theme, first_day, timezone, "values")
        VALUES (?, ?, ?, ?, ?)
        """,
        (
            "client",
            client_values["client.general.theme"],
            client_values["client.general.firstDay"],
            client_values["client.general.timeZone"],
            json.dumps(client_values, sort_keys=True),
        ),
    )


def create_workspace(
    conn: sqlite3.Connection,
    *,
    slug: str,
    name: str,
    description: str,
    now: datetime,
) -> str:
    workspace_id = stable_id("workspace", slug)
    conn.execute(
        "INSERT INTO workspaces (id, name, description) VALUES (?, ?, ?)",
        (workspace_id, name, description),
    )
    conn.execute(
        """
        INSERT INTO workspace_link_rules (workspace_id, position, regexp, link)
        VALUES (?, 0, ?, ?)
        """,
        (workspace_id, r"CTX-(\d+)", "https://example.test/context/$1"),
    )
    return workspace_id


def create_project(
    conn: sqlite3.Connection,
    *,
    slug: str,
    name: str,
    workspace_id: str,
    now: datetime,
    parent_project_id: str | None = None,
) -> str:
    project_id = stable_id("project", slug)
    conn.execute(
        """
        INSERT INTO projects (id, name, parent_id, workspace_id)
        VALUES (?, ?, ?, ?)
        """,
        (project_id, name, parent_project_id, workspace_id),
    )
    return project_id


def create_context(
    conn: sqlite3.Connection,
    *,
    slug: str,
    name: str,
    workspace_id: str,
    description: str,
    now: datetime,
    status: str = "inactive",
    project_id: str | None = None,
) -> str:
    context_id = stable_id("context", slug)
    conn.execute(
        """
        INSERT INTO contexts
          (id, name, workspace_id, status, archived, description, project_id)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        """,
        (context_id, name, workspace_id, status, False, description, project_id),
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
    conn.execute(
        """
        INSERT INTO intervals
          (id, context_id, start, end, duration, status, workspace_id)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        """,
        (
            interval_id,
            context_id,
            sql_time(start),
            sql_time(end),
            duration_ns(duration),
            status,
            workspace_id,
        ),
    )
    return interval_id


def add_context_tag(
    conn: sqlite3.Connection,
    *,
    context_id: str,
    workspace_id: str,
    name: str,
) -> None:
    tag_id = stable_id("tag", workspace_id, name)
    conn.execute("INSERT OR IGNORE INTO tag (id, name) VALUES (?, ?)", (tag_id, name))
    conn.execute(
        "INSERT OR IGNORE INTO context_tags (context_id, tag_id) VALUES (?, ?)",
        (context_id, tag_id),
    )


def create_interval_record(
    conn: sqlite3.Connection,
    *,
    slug: str,
    context_id: str,
    workspace_id: str,
    start: datetime,
    end: datetime | None,
    now: datetime,
    status: str,
) -> str:
    interval_id = stable_id("interval", slug)
    duration = end - start if end is not None and end > start else timedelta(0)
    conn.execute(
        """
        INSERT INTO intervals
          (id, context_id, start, end, duration, status, workspace_id)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        """,
        (
            interval_id,
            context_id,
            sql_time(start),
            sql_time(end) if end is not None else None,
            duration_ns(duration),
            status,
            workspace_id,
        ),
    )
    return interval_id


def seed_project_hierarchy(
    conn: sqlite3.Connection,
    *,
    workspace_id: str,
    slug_prefix: str,
    projects: list[tuple[str, str, str | None]],
    contexts: list[tuple[str, str, str]],
    interval_hour: int,
    now: datetime,
) -> None:
    project_ids: dict[str, str] = {}
    project_names: dict[str, str] = {}
    for project_slug, project_name, parent_slug in projects:
        project_ids[project_slug] = create_project(
            conn,
            slug=f"{slug_prefix}-{project_slug}",
            name=project_name,
            workspace_id=workspace_id,
            parent_project_id=project_ids[parent_slug] if parent_slug is not None else None,
            now=now,
        )
        project_names[project_slug] = project_name

    for index, (project_slug, context_slug, context_name) in enumerate(contexts):
        context_id = create_context(
            conn,
            slug=f"{slug_prefix}-project-{context_slug}",
            name=context_name,
            workspace_id=workspace_id,
            project_id=project_ids[project_slug],
            description=f"Sample context assigned to the {project_names[project_slug]} project.",
            now=now,
        )
        create_interval(
            conn,
            slug=f"{slug_prefix}-project-{context_slug}-interval",
            context_id=context_id,
            workspace_id=workspace_id,
            start=recent_day_start(now, index % SEEDED_DAYS, interval_hour),
            duration=timedelta(minutes=25 + (index % 3) * 5),
            now=now,
        )


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
        if context_index == 0:
            add_context_tag(
                conn,
                context_id=context_id,
                workspace_id=workspace_id,
                name="important",
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

    seed_project_hierarchy(
        conn,
        workspace_id=workspace_id,
        slug_prefix="large",
        projects=[
            ("atlas", "Atlas Platform", None),
            ("atlas-web", "Web Experience", "atlas"),
            ("atlas-web-accessibility", "Accessibility Refresh", "atlas-web"),
            ("atlas-api", "Public API", "atlas"),
            ("operations", "Operations Suite", None),
            ("operations-reporting", "Reporting Dashboard", "operations"),
        ],
        contexts=[
            ("atlas", "architecture", "Platform Architecture"),
            ("atlas-web", "design-system", "Design System"),
            ("atlas-web-accessibility", "keyboard-navigation", "Keyboard Navigation"),
            ("atlas-api", "api-contract", "API Contract"),
            ("operations", "operations-planning", "Operations Planning"),
            ("operations-reporting", "metrics-dashboard", "Metrics Dashboard"),
        ],
        interval_hour=6,
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

    seed_project_hierarchy(
        conn,
        workspace_id=workspace_id,
        slug_prefix="small",
        projects=[
            ("website", "Personal Website", None),
            ("website-content", "Content Refresh", "website"),
            ("website-content-portfolio", "Portfolio Case Studies", "website-content"),
            ("website-hosting", "Hosting Migration", "website"),
            ("home", "Home Operations", None),
            ("home-finance", "Finance Dashboard", "home"),
        ],
        contexts=[
            ("website", "website-planning", "Website Planning"),
            ("website-content", "copywriting", "Copywriting"),
            ("website-content-portfolio", "case-study-review", "Case Study Review"),
            ("website-hosting", "deployment-checklist", "Deployment Checklist"),
            ("home", "household-planning", "Household Planning"),
            ("home-finance", "budget-review", "Budget Review"),
        ],
        interval_hour=7,
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
        slug="broken-interval-missing-context-second",
        context_id="",
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
        slug="broken-interval-missing-context-and-workspace-not-found",
        context_id="",
        workspace_id="missing-workspace-for-interval",
        start=start + timedelta(hours=5),
        duration=timedelta(minutes=18),
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
        SELECT
          first_interval.id,
          second_interval.id,
          first_interval.start,
          first_interval.end,
          second_interval.start,
          second_interval.end
        FROM intervals first_interval
        JOIN intervals second_interval ON first_interval.id < second_interval.id
        WHERE first_interval.start < second_interval.end
          AND second_interval.start < first_interval.end
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

    now = datetime.now(timezone.utc).replace(microsecond=0)
    with sqlite3.connect(output) as conn:
        conn.execute("PRAGMA foreign_keys = ON")
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
        foreign_key_errors = conn.execute("PRAGMA foreign_key_check").fetchall()
        if foreign_key_errors:
            raise RuntimeError(f"Generated database has foreign key errors: {foreign_key_errors}")
        conn.commit()


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Generate a current ctx 0.7.1 relational SQLite database for testing.",
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
