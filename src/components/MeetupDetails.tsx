// src/components/MeetupDetails.tsx
import React from "react";
import {Box, Text} from "ink";
import {Meetup} from "../api";
import {format} from "date-fns";

interface Props {
  meetup: Meetup;
  onBack(): void;
}

const MeetupDetails: React.FC<Props> = ({meetup, onBack}) => (
  <Box flexDirection="column">
    <Text bold>
      {meetup.title} ({format(new Date(meetup.dateTime), "yyyy-MM-dd HH:mm")})
    </Text>
    <Text>Group: {meetup.groupName}</Text>
    <Text>
      Venue: {meetup.venueName}, {meetup.city}, {meetup.state.toUpperCase()}
    </Text>
    <Text>RSVPs: {meetup.rsvpsCount}</Text>
    <Text>
      Link: <Text underline>{meetup.url}</Text>
    </Text>
    <Box marginTop={1} flexDirection="column">
      <Text underline>Description:</Text>
      <Text>{meetup.description}</Text>
    </Box>
    <Box marginTop={1}>
      <Text dimColor>(press “o” to open link, “b” or Esc to go back)</Text>
    </Box>
  </Box>
);

export default MeetupDetails;
