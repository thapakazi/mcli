import React from "react";
import {Box, Text} from "ink";
import {Meetup} from "../../repos/thapakazi/events/cli/src/api";
import {format} from "date-fns";

interface Props {
  meetups: Meetup[];
  selected: number;
}

const MeetupList: React.FC<Props> = ({meetups, selected}) => (
  <Box flexDirection="column" marginBottom={1}>
    {meetups.map((m, i) => (
      <Box key={m.id}>
        <Text color={i === selected ? "cyan" : undefined}>
          {i === selected ? "❯ " : "  "}
        </Text>
        <Text>
          [{m.id.slice(-4)}] {m.title} —{" "}
          {format(new Date(m.dateTime), "yyyy-MM-dd HH:mm")}
        </Text>
      </Box>
    ))}
    {meetups.length === 0 && <Text>(no meetups match)</Text>}
  </Box>
);

export default MeetupList;
