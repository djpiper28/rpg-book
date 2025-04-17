import { H1 } from "@/components/typography/H1";
import { H2 } from "@/components/typography/H2";
import { client } from "@/lib/grpcClient/client";
import { useSettingsStore } from "@/stores/settingsStore";
import { Button, Checkbox } from "@mantine/core";
import { Settings, TriangleAlert } from "lucide-react";
import { useState, useEffect } from "react";

export function SettingsPage() {
  const { settings, setSettings } = useSettingsStore((s) => s);
  const [dirtySettings, setDirtySettings] = useState(settings);
  useEffect(() => {
    setDirtySettings(settings);
  }, [settings]);

  return (
    <>
      <div className="flex flex-row gap-2 items-center">
        <Settings />
        <H1>Settings</H1>
      </div>

      <Checkbox
        label="Dark Mode"
        description="Whether to use dark mode (checked), or light mode (unchecked) for RPG Book."
        checked={dirtySettings.darkMode}
        onChange={(event) => {
          const s = structuredClone(dirtySettings);
          s.darkMode = event.currentTarget.checked;
          setDirtySettings(s);
        }}
      />

      <div className="flex flex-row justify-between gap-3">
        <Button
          variant="default"
          onClick={() => {
            window.location.href = "/";
          }}
        >
          Cancel
        </Button>
        <Button
          onClick={() => {
            client
              .setSettings(dirtySettings)
              .then(() => {
                setSettings(dirtySettings);
                window.location.href = "/";
              })
              .catch(alert);
          }}
        >
          Save
        </Button>
      </div>

      <div className="flex flex-col gap-3 mt-20 m-2 p-5 border-r-2 border border-red-500 rounded-lg">
        <div className="flex flex-row gap-2 items-center">
          <TriangleAlert />
          <H2>Scary Settings - Best Not To Touch</H2>
        </div>
        <Checkbox
          label="Developer Mode"
          checked={dirtySettings.devMode}
          onChange={(event) => {
            const s = structuredClone(dirtySettings);
            s.devMode = event.currentTarget.checked;
            setDirtySettings(s);
          }}
          description="Do not turn developer mode on unless you know what you are doing."
        />
      </div>
    </>
  );
}
