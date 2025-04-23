// src/components/MeetupDetails.tsx
import React, { useState, useEffect } from "react";
import { Box, Text } from "ink";
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

  // title color: blue for online upcoming, greenBright for physical upcoming
  const titleColor = !isPast
    ? isOnline
      ? "blue"
      : "greenBright"
    : undefined;

  // glow animation for link
  const [glowPos, setGlowPos] = useState(0);
  const link = meetup.url;

  useEffect(() => {
    const id = setInterval(() => {
      setGlowPos(pos => (pos + 1) % link.length);
    }, 100);
    return () => clearInterval(id);
  }, [link]);

  const renderGlowingLink = () => (
    <Text>
      {link.split("").map((char, idx) => (
        <Text
          key={idx}
          color={idx === glowPos ? "white" : "gray"}
          bold={idx === glowPos}
        >
          {char}
        </Text>
      ))}
    </Text>
  );

  return (
    <Box flexDirection="column">
      {/* Title */}
      <Text bold color={titleColor} dimColor={isPast}>
        {meetup.title} (
        {format(new Date(meetup.dateTime), "yyyy-MM-dd HH:mm")})
      </Text>

      {/* Group, Venue, RSVPs */}
      <Text color="cyan">{`Group: ${meetup.groupName}`}</Text>
      <Text color="magenta">{`Venue: ${meetup.venueName}, ${meetup.city}, ${meetup.state.toUpperCase()}`}</Text>
      <Text color="yellow">{`RSVPs: ${meetup.rsvpsCount}`}</Text>

      {/* Link with glowing animation */}
      <Box marginTop={1} flexDirection="column">
        <Text>Link:</Text>
        <Box marginLeft={2}>{renderGlowingLink()}</Box>
      </Box>

      {/* Description with spacing */}
      <Box marginTop={1} flexDirection="column" marginLeft={2}>
        {meetup.description.split("\n").map((line, i) => (
          <Text key={i}>{line}</Text>
        ))}
     </Box>
    </Box>
  );
};

export default MeetupDetails;
