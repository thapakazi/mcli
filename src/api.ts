export interface Meetup {
  id: string;
  groupName: string;
  title: string;
  url: string;
  eventType: string;
  dateTime: string;
  description: string;
  imageUrl: string;
  venueName: string;
  venueAddress: string;
  city: string;
  state: string;
  country: string;
  rsvpsCount: number;
  ticketCount: number;
}

const BASE = "http://localhost:3000";

export async function fetchMeetups(): Promise<Meetup[]> {
  const res = await fetch(`${BASE}/meetups`);
  if (!res.ok) throw new Error(res.statusText);
  return res.json();
}

export async function fetchLumas(): Promise<Meetup[]> {
  const res = await fetch(`${BASE}/lumas`);
  if (!res.ok) throw new Error(res.statusText);
  return res.json();
}

export async function fetchEvents(): Promise<Meetup[]> {
  const res = await fetch(`${BASE}/events`);
  if (!res.ok) throw new Error(res.statusText);
  return res.json();
}

export async function fetchMeetupById(id: string): Promise<Meetup> {
  const res = await fetch(`${BASE}/meetup/${id}`);
  if (!res.ok) throw new Error(res.statusText);
  return res.json();
}

export async function fetchLumaById(id: string): Promise<Meetup> {
  const res = await fetch(`${BASE}/luma/${id}`);
  if (!res.ok) throw new Error(res.statusText);
  return res.json();
}

export async function refreshLumaById(id: string): Promise<Meetup> {
  const res = await fetch(`${BASE}/fetch/luma/${id}`);
  if (!res.ok) throw new Error(res.statusText);
  return res.json();
}



export async function fetchMeetupByLocation(location: string): Promise<Meetup> {
  const res = await fetch(`${BASE}/fetch?location=${location}`);
  if (!res.ok) throw new Error(res.statusText);
  return res.json();
}
