import React, { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { X, Plus } from "lucide-react";

export interface EditableBadgesProperties {
    onChange: (badges: string[]) => void
    initBadges: string[]
}

export default function EditableBadges({ onChange, initBadges }: EditableBadgesProperties) {
    const [badges, setBadges] = useState(initBadges ?? []);
    const [newBadge, setNewBadge] = useState("");

    const handleAdd = () => {
        const trimmed = newBadge.trim();
        if (trimmed && !badges.includes(trimmed)) {
            setBadges([...badges, trimmed]);
            onChange([...badges, trimmed])
            setNewBadge("");
        }
    };

    const handleRemove = (badge: string) => {
        const newBadges = badges.filter((b) => b !== badge);
        setBadges(newBadges);
        onChange(newBadges)
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") handleAdd();
    };

    return (
        <div className="w-full max-w-md mx-auto space-y-4">
            <div className="flex flex-wrap gap-2">
                {badges.map((badge) => (
                    <div
                        key={badge}
                    >
                        <Badge
                            variant="secondary"
                            className="flex items-center gap-1 pr-1"
                        >
                            {badge}
                            <button
                                onClick={() => handleRemove(badge)}
                                className="ml-1 rounded-full hover:bg-muted p-0.5"
                            >
                                <X className="h-3 w-3" />
                            </button>
                        </Badge>
                    </div>
                ))}
            </div>

            <div className="flex gap-2">
                <Input
                    placeholder="New label..."
                    value={newBadge}
                    onChange={(e) => setNewBadge(e.target.value)}
                    onKeyDown={handleKeyDown}
                />
                <Button onClick={handleAdd}>
                    <Plus className="h-4 w-4 mr-1" /> Add
                </Button>
            </div>
        </div>
    );
}
