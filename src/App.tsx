// src/App.tsx
import React, {useState, useEffect, useMemo} from "react";
import {Box, useInput, Text, useStdout} from "ink";
import {spawn} from "child_process";
import {fetchMeetups, fetchMeetupById, Meetup} from "./api";
import MeetupList from "./components/MeetupList";
import SearchBar from "./components/SearchBar";
import MeetupDetails from "./components/MeetupDetails";

// hook to track terminal size on resize
function useStdoutDimensions(): [number, number] {
  const getSize = (): [number, number] => [
    process.stdout.columns || 0,
    process.stdout.rows    || 0
  ];
  const [size, setSize] = useState<[number, number]>(getSize());
  useEffect(() => {
    const onResize = () => setSize(getSize());
    process.on("SIGWINCH", onResize);
    return () => void process.off("SIGWINCH", onResize);
  }, []);
  return size;
}

// cross-platform browser opener
function openUrl(url: string) {
  const cmd =
    process.platform === "darwin" ? "open" :
    process.platform === "win32"  ? "start" :
    "xdg-open";
  spawn(cmd, [url], {stdio: "ignore", detached: true}).unref();
}

type View = "list" | "details";

const App: React.FC = () => {
  // — get Ink’s stdout so we can clear the screen —
  const {stdout} = useStdout();

  // — dynamic pageSize = 80% of terminal height —
  const [, rows] = useStdoutDimensions();
  const pageSize = Math.max(1, Math.floor(rows * 0.95));

  // — state —
  const [view, setView]               = useState<View>("list");
  const [meetups, setMeetups]         = useState<Meetup[]>([]);
  const [selected, setSelected]       = useState(0);
  const [offset, setOffset]           = useState(0);
  const [search, setSearch]           = useState("");
  const [isSearching, setIsSearching] = useState(false);
  const [detail, setDetail]           = useState<Meetup | null>(null);
  const [detailsOffset, setDetailsOffset] = useState(0);

  // — CLEAR SCREEN ANY TIME view CHANGES —
  useEffect(() => {
    // CSI 2J = clear screen, CSI 0;0f = move cursor to top left
    stdout.write("\x1B[2J\x1B[0;0f");
  }, [view]);

  // — initial fetch —
  useEffect(() => {
    fetchMeetups().then(setMeetups).catch(console.error);
  }, []);

  // — filter + sort ascending —
  const filtered = useMemo(() => {
    const term = search.toLowerCase();
    return meetups
      .filter(m =>
        m.title.toLowerCase().includes(term) ||
        m.groupName.toLowerCase().includes(term) ||
        m.city.toLowerCase().includes(term)
             )
      .sort((a, b) =>
        new Date(a.dateTime).getTime() - new Date(b.dateTime).getTime()
           );
  }, [meetups, search]);

  // — reset on search or resize —
  useEffect(() => {
    setSelected(0);
    setOffset(0);
  }, [search, pageSize]);

  // — skip past events on list load —
  useEffect(() => {
    if (view !== "list") return;
    const now = Date.now();
    const idx = filtered.findIndex(m => new Date(m.dateTime).getTime() >= now);
    if (idx > 0) {
      setSelected(idx);
      setOffset(idx);
    }
  }, [filtered, view]);

  // — slide window in list —
  useEffect(() => {
    if (selected < offset) {
      setOffset(selected);
    } else if (selected >= offset + pageSize) {
      setOffset(selected - pageSize + 1);
    }
  }, [selected, offset, pageSize]);

  // — reset details scroll when opening —
  useEffect(() => {
    if (detail) setDetailsOffset(0);
  }, [detail]);

  // — LIST view key handling —
  useInput((input, key) => {
    if (view !== "list") return;

    if (key.escape) {
      setIsSearching(prev => !prev);
      return;
    }
    if (input === "/" && !isSearching) {
      setIsSearching(true);
      return;
    }
    if (key.return && isSearching) {
      setIsSearching(false);
      return;
    }
    if (!isSearching && input === "r") {
      fetchMeetups()
        .then(data => {
          setMeetups(data);
          setSelected(0);
          setOffset(0);
        })
        .catch(console.error);
      return;
    }
    if (!isSearching) {
      if (key.upArrow || input === "k") {
        setSelected(i => Math.max(0, i - 1));
      } else if (key.downArrow || input === "j") {
        setSelected(i => Math.min(filtered.length - 1, i + 1));
      } else if (key.return && filtered[selected]) {
        fetchMeetupById(filtered[selected].id).then(m => {
          setDetail(m);
          setView("details");
        });
      }
    }
  });

  // — DETAILS view key handling (scroll + open + back) —
  const totalDetailLines = detail
    ? 5 + 1 + detail.description.split("\n").length + 1
    : 0;
  const maxDetailsOffset = Math.max(0, totalDetailLines - pageSize);

  useInput((input, key) => {
    if (view !== "details") return;

    if (input === "o" && detail) {
      openUrl(detail.url);
    } else if (key.escape || input === "b") {
      setView("list");
    } else if (key.upArrow || input === "k") {
      setDetailsOffset(o => Math.max(0, o - 1));
    } else if (key.downArrow || input === "j") {
      setDetailsOffset(o => Math.min(maxDetailsOffset, o + 1));
    }
  });

  // — render —
  const visibleList = filtered.slice(offset, offset + pageSize);
  const selectedInWindow = selected - offset;

  return (
    <Box flexDirection="column">
      {view === "list" && (
        <>
          <MeetupList
        filtered={visibleList}
        totalCount={pageSize}
        selected={selectedInWindow}
          />
          <SearchBar
        value={search}
        onChange={setSearch}
        placeholder="Filter…"
        focus={isSearching}
        onSubmit={() => setIsSearching(false)}
          />
          <Box marginTop={1}>
          <Text dimColor>(r refresh, / search, j/k or ↑/↓ move)</Text>
          </Box>
          </>
      )}

    {view === "details" && detail && (
      <>
        <MeetupDetails
      meetup={detail}
      offset={detailsOffset}
      pageSize={pageSize}
        />
        {/* Footer hint */}
        <Box marginTop={1}>
        <Text dimColor>
        (o open link, j/k or ↑/↓ scroll, b or Esc to go back)
      </Text>
        </Box>
        </>
    )}
    </Box>
  );
};

export default App;
