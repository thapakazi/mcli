import React, {useEffect, useState, useMemo} from "react";
import {Box, Text, useInput, render} from "ink";
import {fetchMeetups, fetchMeetupById, Meetup} from "./api";
import MeetupList from "../../../../../.emacs.d/backups/!Users!thapakazi!repos!thapakazi!events!cli!src!components!MeetupList.tsx~";
import SearchBar from "./components/SearchBar";
import MeetupDetails from "../../../../../.emacs.d/backups/!Users!thapakazi!repos!thapakazi!events!cli!src!components!MeetupDetails.tsx~";

type View = "list" | "details";

const App: React.FC = () => {
  const [view, setView] = useState<View>("list");
  const [meetups, setMeetups] = useState<Meetup[]>([]);
  const [selected, setSelected] = useState(0);
  const [search, setSearch] = useState("");
  const [detail, setDetail] = useState<Meetup | null>(null);

  // fetch list once
  useEffect(() => {
    fetchMeetups().then(setMeetups).catch(err => {
      console.error(err);
      process.exit(1);
    });
  }, []);

  // filtered list
  const filtered = useMemo(() => {
    const term = search.toLowerCase();
    return meetups.filter(m =>
      m.title.toLowerCase().includes(term) ||
      m.groupName.toLowerCase().includes(term) ||
      m.city.toLowerCase().includes(term)
    );
  }, [meetups, search]);

  // handle keys
  useInput((input, key) => {
    if (view === "list") {
      if (key.upArrow) {
        setSelected(i => Math.max(0, i - 1));
      } else if (key.downArrow) {
        setSelected(i => Math.min(filtered.length - 1, i + 1));
      } else if (key.return) {
        const m = filtered[selected];
        fetchMeetupById(m.id).then(m => {
          setDetail(m);
          setView("details");
        });
      }
    } else if (view === "details") {
      if (key.escape || input === "b") {
        setView("list");
      }
    }
  });

  return (
    <Box flexDirection="column">
      {view === "list" && (
        <>
          <MeetupList meetups={filtered} selected={selected} />
          <SearchBar value={search} onChange={setSearch} placeholder="Filter…" />
        </>
      )}
      {view === "details" && detail && (
        <MeetupDetails meetup={detail} onBack={() => setView("list")} />
      )}
    </Box>
  );
};

export default App;
