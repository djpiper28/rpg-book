const nouns = [
  "Vampires",
  "Werewolves",
  "Mummies",
  "Goblins",
  "Wizards",
  "Adventurers",
  "Priests",
  "Fools",
  "Rascals",
  "Mercenaries",
  "Dwarves",
  "Grots",
  "Dire-Badgers",
  "Servants",
  "Clerics",
  "Barbarians",
  "Warlords",
  "Peasants",
  "Serfs",
  "Nomads",
  "Farmers",
  "Clergy",
  "Userpers",
  "Warriors",
  "Mages",
  "Knights",
  "Crooks",
  "Criminals",
  "Outcasts",
  "Students",
  "Worshippers",
  "Dragons",
  "Kobolds",
  "Mormons",
  "Punters",
  "Thieves",
  "Rogues",
  "Morons",
];

const ingVerbs = ["Awakening", "Reckoning", "Masquerade", "Happening"];

const places = [
  "Waterdeep",
  "Siberia",
  "Yorkshire",
  "The Gaslands",
  "The City",
  "The Dark Valley",
  "The Light Valley",
  "The River Bed",
  "The Evil Island",
  "The Twee Village",
  "The Model Village",
  "The Dungeons",
  "The Wizard Tower",
  "The Wizard Academy",
  "The Ocean",
  "Distant Lands",
  "Hell",
  "Heaven",
  "The Desert",
  "The Coast",
  "Savannah",
  "Synnibarr",
];

const adjectives = [
  "Poverty",
  "Fortune",
  "Holyness",
  "Valour",
  "Noble-Birth",
  "Lesser-Birth",
  "Harmony",
  "Peace",
  "War",
  "Cowardice",
  "Bravery",
  "Worth",
  "Ruin",
  "Riches",
  "Desperation",
  "Worm",
  "Blood",
  "Courage",
];

function random<T>(x: T[]) {
  const el = Math.round(Math.random() * (x.length - 1));
  return x[el];
}

function aTheBGenerator(): string {
  return `${random(nouns)} The ${random(ingVerbs)}`;
}

function aOfBGenerator(): string {
  return `${random(nouns)} of ${random([...places, ...adjectives])}`;
}

function aAndBGenerator(): string {
  return `${random(nouns)} And ${random(nouns)}`;
}

function aGenerator(): string {
  return random(nouns);
}

const generators = [aTheBGenerator, aOfBGenerator, aAndBGenerator, aGenerator];

export function randomPlace(): string {
  return random(generators)();
}
