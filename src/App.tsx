// src/App.tsx
import React, {useState, useEffect, useMemo} from "react";
import {Box, useInput, Text, useStdout} from "ink";
import TextInput from "ink-text-input";
import {spawn} from "child_process";
import {format} from "date-fns";
import {
  fetchEvents,
  fetchMeetups,
  fetchLumas,
  fetchMeetupById,
  fetchLumaById,
  fetchMeetupByLocation,
  refreshLumaById,
  Meetup
} from "./api";
import Layout from "./components/Layout";
import SearchBar from "./components/SearchBar";
import MeetupList from "./components/MeetupList";
import MeetupDetails from "./components/MeetupDetails";

type View = "list" | "details";
type SearchMode = "filter" | "fetch";

const App: React.FC = () => {
  // ─── Core state ─────────────────────────────────────────────────────────────
  const [view, setView]                 = useState<View>("list");
  const [meetups, setMeetups]           = useState<Meetup[]>([]);
  const [selected, setSelected]         = useState(0);
  const [offset, setOffset]             = useState(0);
  const [detail, setDetail]             = useState<Meetup | null>(null);
  const [detailsOffset, setDetailsOffset] = useState(0);

  // filter vs fetch mode
  const [mode, setMode]                 = useState<SearchMode>("filter");
  const [isFocused, setIsFocused]       = useState(false);
  const [search, setSearch]             = useState("");
  const [formLocation, setFormLocation] = useState("");

  // ─── Paging setup ────────────────────────────────────────────────────────────
  const {stdout} = useStdout();
  const rows      = (stdout.rows ?? 24) - 6; // leave room for header/footer
  const pageSize  = Math.max(1, Math.floor(rows));

  // ─── Initial fetch ───────────────────────────────────────────────────────────
  useEffect(() => {
    fetchEvents().then(setMeetups).catch(console.error);
  }, []);

  // ─── Filter + sort ───────────────────────────────────────────────────────────
  const filtered = useMemo(() => {
    const term = search.toLowerCase();
    return meetups
      .filter(m =>
        m.title?.toLowerCase().includes(term) ||
        m.groupName?.toLowerCase().includes(term) ||
        m.city?.toLowerCase().includes(term) ||
        m.venueCity?.toLowerCase().includes(term)
      )
      .sort((a, b) =>
        new Date(a.dateTime).getTime() - new Date(b.dateTime).getTime()
      );
  }, [meetups, search]);

  // ─── Reset list cursor/window on search/mode/pageSize change ────────────────
  useEffect(() => {
    setSelected(0);
    setOffset(0);
  }, [search, mode, pageSize]);

  // ─── Skip past events on first list entry ────────────────────────────────────
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

  // ─── Slide window in list ────────────────────────────────────────────────────
  useEffect(() => {
    if (selected < offset) {
      setOffset(selected);
    } else if (selected >= offset + pageSize) {
      setOffset(selected - pageSize + 1);
    }
  }, [selected, offset, pageSize]);

  // ─── Reset details scroll when opening a meetup ──────────────────────────────
  useEffect(() => {
    if (detail) setDetailsOffset(0);
  }, [detail]);

  // ─── LIST + FORM input ───────────────────────────────────────────────────────
  useInput((input, key) => {
    if (view !== "list") return;

    // 1) Toggle focus: Esc blurs/unfocuses
    if (key.escape && isFocused) {
      setIsFocused(false);
      return;
    }

    // 2) "/" enters filter mode
    if (!isFocused && input === "/") {
      setMode("filter");
      setIsFocused(true);
      setSearch("");
      return;
    }

    // 3) "f" enters fetch mode
    if (!isFocused && input === "f") {
      setMode("fetch");
      setIsFocused(true);
      setFormLocation("");
      return;
    }

    // 4) if focused, all keys go to TextInput:
    if (isFocused) return;

    // 5) refresh on "r"
    if (input === "r") {
      fetchEvents()
        .then(data => {
          setMeetups(data);
          setDetail(null);
          setView("list");
        })
        .catch(console.error);
      return;
    }

    // 6) navigation & dive-in
    if (key.upArrow || input === "k") {
      setSelected(i => Math.max(0, i - 1));
    } else if (key.downArrow || input === "j") {
      setSelected(i => Math.min(filtered.length - 1, i + 1));
    } else if (key.return && filtered[selected]) {
      let source =  filtered[selected].source;
      if ( source == 'luma') {
        fetchLumaById(filtered[selected].id).then(m => {
          setDetail(m);
          setView("details");
        });
      } else {
        fetchMeetupById(filtered[selected].id).then(m => {
          setDetail(m);
          setView("details");
        });
      }
    }
  });

  // ─── SUBMIT handler for unified SearchBar ────────────────────────────────────
  const handleSubmit = async (value: string) => {
    setIsFocused(false);

    if (mode === "filter") {
      // just leave `search` as-is
      return;
    } else {
      // fetch by location
      try {
        const meetup = await fetchMeetupByLocation(value);
        setMeetups([meetup]);
      } catch (err) {
        console.error("Fetch by location failed:", err);
      }
    }
  };

  // ─── DETAILS view input (scroll, open, back) ────────────────────────────────
  const totalDetailLines = detail
    ? 5 + 1 + detail.description?.split("\n").length + 1
    : 0;
  const maxDetailsOffset = Math.max(0, totalDetailLines - pageSize);

  useInput((input, key) => {
    if (view !== "details" || !detail) return;

    // “r” reloads from Luma
    if (input === "r") {
      if (detail.source == "luma" && detail.description === null) {
        refreshLumaById(detail.id)
          .then(m => {
            debugger;
            return setDetail(m)
          })
          .catch(err => console.error("Luma refresh failed:", err));
        return;
      }
    }

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

  // ─── Render ─────────────────────────────────────────────────────────────────
  const visibleList     = filtered.slice(offset, offset + pageSize);
  const selectedInWindow = selected - offset;

  const header = (
    <Text bold>
      {view === "list" ? "📅 Meetups" : "🔎 Meetup Details"}
    </Text>
  );

  // unified SearchBar in footer
  const footer = (
    <Box
      flexDirection="row"
      justifyContent="space-between"
      borderColor="magenta"
      borderStyle="round"
      width="100%"
    >
      <Box flexGrow={1}>
        <SearchBar
          value={mode === "filter" ? search : formLocation}
          onChange={mode === "filter" ? setSearch : setFormLocation}
          placeholder={
            mode === "filter" ? "Filter meetups…" : "Fetch location…"
          }
          focus={isFocused}
          onSubmit={handleSubmit}
        />
      </Box>
      <Text dimColor>
        {format(new Date(), "yyyy-MM-dd HH:mm:ss")}
      </Text>
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
