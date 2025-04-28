// src/components/MeettsxupList.tsx
import React from "react";
import { Box, Text } from "ink";
import gradient from 'gradient-string';
import { Meetup } from "../api";
import { format, differenceInCalendarDays, differenceInHours } from "date-fns";

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

        let timeToGo: string;
        const eventDate = new Date(m.dateTime);
        if (eventDate.getTime() <= now) {
          timeToGo = "started";
        } else {
          const daysDiff = differenceInCalendarDays(eventDate, now);

          if (daysDiff === 0) {
            // same calendar day
            timeToGo = `today @${format(eventDate, "HH:mm")}`;
          } else {
            // full days + remaining hours
            const totalHours = differenceInHours(eventDate, now);
            const hoursOnly = totalHours - daysDiff * 24;
            timeToGo = `${daysDiff}d`;
          }
        }

        // decide color: blue for online upcoming, bright green for physical upcoming
        const color = !isPast
          ? isOnline
            ? "blue"
            : "green"
          : undefined;

        const source = m.source == "luma" ? gradient(['#6e2fe3','#0cabf7','#e27417','#1f6f05'])('✦︎') : gradient(['#f6405fcf','pink'])('☘︎');
        return (
          <Box key={m.id}>
            {/* arrow indicator */}
            <Text color={i === selected ? "cyan" : undefined}>
              {i === selected ? "❯ " : "  "}
            </Text>

            {/* date + title with conditional coloring/dimming */}
            <Text color={color} dimColor={isPast}>
            {timeToGo} {source} {m.title} | {m.venueName}
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
