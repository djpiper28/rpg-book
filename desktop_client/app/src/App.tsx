import { Settings } from "lucide-react";
import { client, logger } from "./lib/grpcClient/client";
import { Button, MantineProvider, Table, Title } from "@mantine/core";
import { P } from "./components/typography/P";
import { H2 } from "./components/typography/H2";
import { useEffect } from "react";
import { useSettingsStore } from "./stores/settingsStore";

function App() {
  const { settings, setSettings } = useSettingsStore((s) => s);
  useEffect(() => {
    client
      .getSettings({})
      .then((x) => {
        setSettings(x.response);
      })
      .catch((e) => {
        logger.error(`Cannot get settings ${e}`, {});
      });
  }, [setSettings]);

  return (
    <MantineProvider
      withCSSVariables
      withGlobalStyles
      withNormalizeCSS
      defaultColorScheme={settings.darkMode ? "dark" : "light"}
    >
      <div className="flex flex-col gap-3 p-2">
        <div
          id="menu"
          className="flex flex-row gap-3 justify-between items-center"
        >
          <Title>RPG Book</Title>
          <Button aria-label="Open settings" variant="filled">
            <Settings /> Settings
          </Button>
        </div>

        <H2>Recent Projects:</H2>
        <P>Lorem ipsum sit amet dolor</P>
        <Table variant="vertical">
          <Table.Thead>
            <Table.Tr>
              <Table.Th>Project Name</Table.Th>
              <Table.Th>Last Opened</Table.Th>
              <Table.Th>Size</Table.Th>
              <Table.Th>Location</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            <Table.Tr>
              <Table.Th>Test Project</Table.Th>
              <Table.Th>3 days ago</Table.Th>
              <Table.Th>5 MiB</Table.Th>
              <Table.Th>/home/example/test.rpg</Table.Th>
            </Table.Tr>
          </Table.Tbody>
        </Table>
      </div>
    </MantineProvider>
  );
}

export default App;
