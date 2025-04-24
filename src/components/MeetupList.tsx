// src/components/MeettsxupList.tsx
import React from "react";
import { Box, Text } from "ink";
import { Meetup } from "../api";
import { format } from "date-fns";

interface Props {
  filtered: Meetup[];
  totalCount: number;
  selected: number;
}

const MeetupList: React.FC<Props> = ({ filtered, totalCount, selected }) => {
  const now = Date.now();

  return (
    <Box flexDirection="column" height={totalCount} >
      {filtered.map((m, i) => {
        const eventTime = new Date(m.dateTime).getTime();
        const isPast    = eventTime < now;
        const isOnline  = m.venueName === "Online event";

        // decide color: blue for online upcoming, bright green for physical upcoming
        const color = !isPast
          ? isOnline
            ? "blue"
            : "greenBright"
          : undefined;

        return (
          <Box key={m.id}>
            {/* arrow indicator */}
            <Text color={i === selected ? "cyan" : undefined}>
              {i === selected ? "❯ " : "  "}
            </Text>

            {/* date + title with conditional coloring/dimming */}
            <Text color={color} dimColor={isPast}>
              {format(new Date(m.dateTime), "yyyy-MM-dd HH:mm")} – {m.title} | {m.rsvpsCount} - {m.venueName}, {m.city}
            </Text>
          </Box>
        );
      })}

      {/* pad with blanks so Ink erases old rows */}
      {Array.from({ length: totalCount - filtered.length }).map((_, idx) => (
        <Text key={`empty-${idx}`}> </Text>
      ))}
    </Box>
  );
};

export default MeetupList;
