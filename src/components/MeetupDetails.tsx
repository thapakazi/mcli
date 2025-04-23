import React from "react";
import {Box, Text} from "ink";
import {Meetup} from "../api";
import {format} from "date-fns";

interface Props {
  meetup: Meetup;
  offset: number;
  pageSize: number;
}

const ScrollableMeetupDetails: React.FC<Props> = ({meetup, offset, pageSize}) => {
  // flatten the meetup into a list of lines
  const lines: string[] = [];
  lines.push(
    `${meetup.title} (${format(new Date(meetup.dateTime), "yyyy-MM-dd HH:mm")})`
  );
  lines.push(`Group: ${meetup.groupName}`);
  lines.push(
    `Venue: ${meetup.venueName}, ${meetup.city}, ${meetup.state.toUpperCase()}`
  );
  lines.push(`RSVPs: ${meetup.rsvpsCount}`);
  lines.push(`Link: ${meetup.url}`);
  lines.push("");
  lines.push("Description:");
  lines.push(...meetup.description.split("\n"));

  // slice out the window
  const slice = lines.slice(offset, offset + pageSize);

  return (
    <Box flexDirection="column">
      {slice.map((line, i) => (
        <Text key={i}>{line}</Text>
      ))}

      {/* pad with blank lines so we always render exactly pageSize rows */}
      {Array.from({ length: pageSize - slice.length }).map((_, i) => (
        <Text key={`pad-${i}`}> </Text>
      ))}

      <Box marginTop={1}>
        <Text dimColor>
          (j/k or ↑/↓ to scroll, “o” to open link, “b” or Esc to go back)
        </Text>
      </Box>
    </Box>
  );
};

export default ScrollableMeetupDetails;
