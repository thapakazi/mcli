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

type View = "list" | "details";

const App: React.FC = () => {
  // Core state
  const [view, setView] = useState<View>("list");
  const [meetups, setMeetups]             = useState<Meetup[]>([]);
  const [selected, setSelected]           = useState(0);
  const [offset, setOffset]               = useState(0);
  const [search, setSearch]               = useState("");
  const [isSearching, setIsSearching]     = useState(false);
  const [detail, setDetail]               = useState<Meetup | null>(null);
  const [detailsOffset, setDetailsOffset] = useState(0);

  const { stdout } = useStdout();
  const rows = stdout.rows - 10 ?? 24;
  const pageSize = Math.max(1, Math.floor(rows));

  // Fetch once
  useEffect(() => {
    fetchMeetups().then(setMeetups).catch(console.error);
  }, []);

  // Filter + sort
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

  // Reset on search or resize-derived pageSize
  useEffect(() => {
    setSelected(0);
    setOffset(0);
  }, [search, pageSize]);

  // Skip past events on first entering list
  useEffect(() => {
    if (view !== "list") return;
    const now = Date.now();
    const idx = filtered.findIndex(
      m => new Date(m.dateTime).getTime() >= now
    );
    if (idx > 0) {
      setSelected(idx);
      setOffset(idx);
    }
  }, [filtered, view]);

  // Slide window for list
  useEffect(() => {
    if (selected < offset) {
      setOffset(selected);
    } else if (selected >= offset + pageSize) {
      setOffset(selected - pageSize + 1);
    }
  }, [selected, offset, pageSize]);

  //  Reset details scroll on open
  useEffect(() => {
    if (detail) setDetailsOffset(0);
  }, [detail]);

  // LIST view input
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
      fetchMeetups().then(data => {
        setMeetups(data);
        setSelected(0);
        setOffset(0);
      }).catch(console.error);
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

  // 12) DETAILS view input (scroll + open + back)
  const totalDetailLines = detail
    ? 5 + 1 + detail.description.split("\n").length + 1
    : 0;
  const maxDetailsOffset = Math.max(0, totalDetailLines - pageSize);

  useInput((input, key) => {
    if (view !== "details" || !detail) return;

    if (input === "o") {
      const cmd =
        process.platform === "darwin" ? "open" :
        process.platform === "win32"  ? "start" :
                                        "xdg-open";
      spawn(cmd, [detail.url], {stdio:"ignore", detached:true}).unref();
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
  );

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
      <Box>
        <Text dimColor>
          {format(new Date(), "yyyy-MM-dd HH:mm:ss")}
        </Text>
      </Box>
    </Box>
  );

  return (
    <Layout header={header} footer={footer}>
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
