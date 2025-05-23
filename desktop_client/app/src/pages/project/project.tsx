import { Tabs } from "@mantine/core";
import { useNavigate } from "react-router";
import { H2 } from "@/components/typography/H2";
import { useTabStore } from "@/stores/tabStore";
import { CharacterTab } from "./tabs/characterTab";

export function ProjectPage() {
  const navigate = useNavigate();
  const tabs = useTabStore();

  if (!tabs.selectedTab) {
    // eslint-disable-next-line @typescript-eslint/no-floating-promises
    navigate("/");
    return;
  }

  const currentTab = tabs.tabs[tabs.selectedTab.id];

  return (
    <>
      <Tabs defaultValue="characters">
        <Tabs.List>
          <H2>{currentTab.name}</H2>
          <Tabs.Tab value="notes">Notes</Tabs.Tab>
          <Tabs.Tab value="characters">Characters</Tabs.Tab>
          <Tabs.Tab value="items">Items</Tabs.Tab>
          <Tabs.Tab value="factions">Factions</Tabs.Tab>
          <Tabs.Tab value="maps">Maps &amp; Places</Tabs.Tab>
          <Tabs.Tab value="settings">Settings</Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="characters">
          <CharacterTab />
        </Tabs.Panel>
      </Tabs>
    </>
  );
}
