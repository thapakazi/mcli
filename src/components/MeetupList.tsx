// src/components/MeetupList.tsx
import React from "react";
import {Box, Text} from "ink";
import {Meetup} from "../api";
import {format} from "date-fns";

interface Props {
  filtered: Meetup[];
  totalCount: number;
  selected: number;
}

const MeetupList: React.FC<Props> = ({ filtered, totalCount, selected }) => {
  const now = Date.now();

  return (
    <Box flexDirection="column" height={totalCount}>
      {filtered.map((m, i) => {
        const eventTime = new Date(m.dateTime).getTime();
        const isPast = eventTime < now;

        return (
          <Box key={m.id}>
            {/* arrow stays highlighted if selected */}
            <Text color={i === selected ? "cyan" : undefined}>
              {i === selected ? "❯ " : "  "}
            </Text>

            {/* dim the whole row if it's in the past */}
            <Text dimColor={isPast}>
              {format(new Date(m.dateTime), "yyyy-MM-dd HH:mm")} – {m.title}
            </Text>
          </Box>
        );
      })}

      {/* pad with blanks so old rows get wiped */}
      {Array.from({ length: totalCount - filtered.length }).map((_, idx) => (
        <Text key={`empty-${idx}`}> </Text>
      ))}
    </Box>
  );
};

export default MeetupList;
