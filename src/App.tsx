// src/App.tsx
import React, {useState, useEffect, useMemo} from "react";
import {Box, useInput, Text} from "ink";
import {spawn} from "child_process";
import {fetchMeetups, fetchMeetupById, Meetup} from "./api";
import MeetupList from "./components/MeetupList";
import SearchBar from "./components/SearchBar";
import MeetupDetails from "./components/MeetupDetails";

// —————————————————————————————————————————————————————————————————————————————
// Custom hook to get terminal [cols, rows] and update on resize
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

// —————————————————————————————————————————————————————————————————————————————
// Cross-platform URL opener
function openUrl(url: string) {
  const cmd =
    process.platform === "darwin" ? "open" :
    process.platform === "win32"  ? "start" :
                                    "xdg-open";
  spawn(cmd, [url], { stdio: "ignore", detached: true }).unref();
}

type View = "list" | "details";

const App: React.FC = () => {
  // 1) Dynamic pageSize = 80% of terminal height
  const [, rows] = useStdoutDimensions();
  const pageSize = Math.max(1, Math.floor(rows * 0.95));

  // 2) State
  const [view, setView]           = useState<View>("list");
  const [meetups, setMeetups]     = useState<Meetup[]>([]);
  const [selected, setSelected]   = useState(0);    // index in `filtered`
  const [offset, setOffset]       = useState(0);    // window start
  const [search, setSearch]       = useState("");
  const [detail, setDetail]       = useState<Meetup | null>(null);

  // 3) Fetch once
  useEffect(() => {
    fetchMeetups()
      .then(setMeetups)
      .catch(err => {
        console.error(err);
        process.exit(1);
      });
  }, []);

  // 4) Filter + sort chronologically
  const filtered = useMemo(() => {
    const term = search.toLowerCase();
    return meetups
      .filter(
        m =>
          m.title.toLowerCase().includes(term) ||
          m.groupName.toLowerCase().includes(term) ||
          m.city.toLowerCase().includes(term)
          // fix me: do exact match
          //m.description.toLowerCase().includes(term)
      )
      .sort(
        (a, b) =>
          new Date(a.dateTime).getTime() - new Date(b.dateTime).getTime()
      );
  }, [meetups, search]);

  // 5) Reset selection/window on search or resize
  useEffect(() => {
    setSelected(0);
    setOffset(0);
  }, [search, pageSize]);

  // 6) **New**: skip all past (dimmed) events
  useEffect(() => {
    if (view !== "list") return;
    const now = Date.now();
    const firstFuture = filtered.findIndex(
      m => new Date(m.dateTime).getTime() >= now
    );
    if (firstFuture > 0) {
      setSelected(firstFuture);
      setOffset(firstFuture);
    }
  }, [filtered, view]);

  // 7) Slide the window if you move out of it
  useEffect(() => {
    if (selected < offset) {
      setOffset(selected);
    } else if (selected >= offset + pageSize) {
      setOffset(selected - pageSize + 1);
    }
  }, [selected, offset, pageSize]);

  // 8) Key handling: ↑/↓/Enter in list, o/Esc/b in details
  useInput((input, key) => {
    if (view === "list") {
      if (key.upArrow) {
        setSelected(i => Math.max(0, i - 1));
      } else if (key.downArrow) {
        setSelected(i => Math.min(filtered.length - 1, i + 1));
      } else if (key.return && filtered[selected]) {
        fetchMeetupById(filtered[selected].id).then(m => {
          setDetail(m);
          setView("details");
        });
      }
    } else {
      if (input === "o" && detail) {
        openUrl(detail.url);
      } else if (key.escape || input === "b") {
        setView("list");
      }
    }
  });

  // 9) Render only the current “page” of items
  const visible = filtered.slice(offset, offset + pageSize);
  const selectedInWindow = selected - offset;

  return (
    <Box flexDirection="column">
      {view === "list" && (
        <>
          <MeetupList
            filtered={visible}
            totalCount={pageSize}
            selected={selectedInWindow}
          />
          <Box marginTop={1}>
            <Text dimColor>
              Showing {offset + 1}–{Math.min(offset + pageSize, filtered.length)} of{" "}
              {filtered.length}
            </Text>
          </Box>
          <SearchBar
            value={search}
            onChange={setSearch}
            placeholder="Filter by title, desc, group or city…"
          />
        </>
      )}

      {view === "details" && detail && (
        <MeetupDetails meetup={detail} onBack={() => setView("list")} />
      )}
    </Box>
  );
};

export default App;
