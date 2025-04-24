// src/App.tsx
import React, {useState, useEffect, useMemo} from "react";
import {Box, useInput, Text, useStdout} from "ink";
import {spawn} from "child_process";
import {format} from "date-fns";
import {fetchMeetups, fetchMeetupById, Meetup} from "./api";
import Layout from "./components/Layout";
import MeetupList from "./components/MeetupList";
import SearchBar from "./components/SearchBar";
import MeetupDetails from "./components/MeetupDetails";

// ─── Helpers ────────────────────────────────────────────────────────────────
// Hook: track terminal [cols, rows] and update on SIGWINCH
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

// Cross‐platform URL opener
function openUrl(url: string) {
  const cmd =
    process.platform === "darwin" ? "open" :
    process.platform === "win32"  ? "start" :
                                    "xdg-open";
  spawn(cmd, [url], {stdio: "ignore", detached: true}).unref();
}

type View = "list" | "details";

const App: React.FC = () => {
  // 1) Grab Ink’s stdout and view state
  const {stdout} = useStdout();
  const [view, setView] = useState<View>("list");

  // 2) Clear screen on view change
  useEffect(() => {
    // ESC[2J clears, ESC[0;0f homes
    stdout.write("\x1B[2J\x1B[0;0f");
  }, [stdout, view]);

  // 3) Compute dynamic pageSize = 80% of terminal height
  const [, rows] = useStdoutDimensions();
  const pageSize = Math.max(1, Math.floor(rows * 0.9));

  // 4) Core state
  const [meetups, setMeetups]         = useState<Meetup[]>([]);
  const [selected, setSelected]       = useState(0);    // index in filtered[]
  const [offset, setOffset]           = useState(0);    // list window start
  const [search, setSearch]           = useState("");
  const [isSearching, setIsSearching] = useState(false);
  const [detail, setDetail]           = useState<Meetup | null>(null);
  const [detailsOffset, setDetailsOffset] = useState(0);

  // 5) Fetch the meetup list once
  useEffect(() => {
    fetchMeetups().then(setMeetups).catch(console.error);
  }, []);

  // 6) Filter + sort chronologically
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

  // 7) Reset list cursor & window on search or resize
  useEffect(() => {
    setSelected(0);
    setOffset(0);
  }, [search, pageSize]);

  // 8) On entering list view, skip past events
  useEffect(() => {
    if (view !== "list") return;
    const now = Date.now();
    const idx = filtered.findIndex(m => new Date(m.dateTime).getTime() >= now);
    if (idx > 0) {
      setSelected(idx);
      setOffset(idx);
    }
  }, [filtered, view]);

  // 9) Slide the list window if the selection moves out of view
  useEffect(() => {
    if (selected < offset) {
      setOffset(selected);
    } else if (selected >= offset + pageSize) {
      setOffset(selected - pageSize + 1);
    }
  }, [selected, offset, pageSize]);

  // 10) Reset details scroll when opening a meetup
  useEffect(() => {
    if (detail) setDetailsOffset(0);
  }, [detail]);

  // 11) LIST view input (j/k, arrows, Enter, r, /, Esc)
  useInput((input, key) => {
    if (view !== "list") return;

    if (key.escape) {
      // toggle search focus
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
      // refresh list
      fetchMeetups().then(data => {
        setMeetups(data);
        setSelected(0);
        setOffset(0);
      }).catch(console.error);
      return;
    }
    if (!isSearching) {
      // navigation & drill-in
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

  // 12) DETAILS view input (scroll j/k, ↑/↓, open, back)
  // compute max scroll offset
  const totalDetailLines = detail
    ? 5 + 1 + detail.description.split("\n").length + 1
    : 0;
  const maxDetailsOffset = Math.max(0, totalDetailLines - pageSize);

  useInput((input, key) => {
    if (view !== "details" || !detail) return;

    if (input === "o") {
      openUrl(detail.url);
    } else if (key.escape || input === "b") {
      setView("list");
    } else if (key.upArrow || input === "k") {
      setDetailsOffset(o => Math.max(0, o - 1));
    } else if (key.downArrow || input === "j") {
      setDetailsOffset(o => Math.min(maxDetailsOffset, o + 1));
    }
  });

  // 13) Render
  const visibleList = filtered.slice(offset, offset + pageSize);
  const selectedInWindow = selected - offset;


  const header = (
    <Text bold>
      {view === "list" ? "📅 Meetups" : "🔎 Meetup Details"}
    </Text>
  )
  
  // Footer: search bar + status bar (current datetime)
  const footer = (
    <Box
      flexDirection="row"
      justifyContent="space-between"
      borderColor="magenta"
      borderStyle="round"
      width="100%"
    >
      <SearchBar
        value={search}
        onChange={setSearch}
        placeholder="Filter meetups…"
        focus={isSearching}
        onSubmit={() => setIsSearching(false)}
      />
      <Box alignItems="flex-end">
        <Text dimColor>
          {format(new Date(), "yyyy-MM-dd HH:mm:ss")}
        </Text>
      </Box>
    </Box>
  );

  return (
    <Layout
      header={header}
      footer={footer}
    >
      {view === "list" && (
        <MeetupList
          filtered={visibleList}
          totalCount={pageSize}
          selected={selectedInWindow}
        />
      )}
      {view === "details" && detail && (
        <MeetupDetails
          meetup={detail}
          offset={detailsOffset}
          pageSize={pageSize}
        />
      )}
    </Layout>
  );
};

export default App;
