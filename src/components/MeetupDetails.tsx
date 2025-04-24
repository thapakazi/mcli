// src/components/MeetupDetails.tsx
import React, { useState, useEffect } from "react";
import { Box, Text, Spacer} from "ink";
import gradient from 'gradient-string';
import { Meetup } from "../api";
import { format } from "date-fns";

interface Props {
  meetup: Meetup;
  onBack(): void;
}

const MeetupDetails: React.FC<Props> = ({ meetup, onBack }) => {
  const now = Date.now();
  const eventTime = new Date(meetup.dateTime).getTime();
  const isPast = eventTime < now;
  const isOnline = meetup.venueName === "Online event";

  const titleColor = !isPast
    ? isOnline ? "blue" : "greenBright"
    : undefined;

  return (
    <Box flexDirection="column" padding={2}>

      {/* Head */}
      <Box flexDirection="column">
        <Text bold color={titleColor} dimColor={isPast}>{meetup.title}</Text>
        <Text color="magenta">📍 {meetup.venueName}, {meetup.city}, {meetup.state.toUpperCase()}</Text>
        <Text color="cyan">📆 {format(new Date(meetup.dateTime), "yyyy-MM-dd || HH:mm")}</Text>
        <Text>🤹 {gradient(['red',' yellow'])(meetup.groupName)}</Text>
        <Text color="yellow">＃ {meetup.rsvpsCount}</Text>
        <Text>🔗 {gradient(['cyan',' pink',' magenta',' red',' yellow',' green',' blue'])(meetup.url)}</Text>
      </Box>

      {/* Description */}
      <Box marginTop={1} flexDirection="column">
        {meetup.description.split("\n").map((line, i) => (
          <Text key={i}>{line}</Text>
        ))}
      </Box>
	  <Spacer />

      {/* Tail */}
      <Box marginTop={1}>
        <Text dimColor>(press “b” or Esc to go back)</Text>
      </Box>

    </Box>
  );
};

export default MeetupDetails;
