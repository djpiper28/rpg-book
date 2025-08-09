import { Button, Checkbox, InputDescription } from "@mantine/core";
import { Settings, TriangleAlert } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { H2 } from "@/components/typography/H2";
import { getSystemClient } from "@/lib/grpcClient/client";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useSettingsStore } from "@/stores/settingsStore";

export function Component() {
  const { setSettings, settings } = useSettingsStore((s) => s);
  const [dirtySettings, setDirtySettings] = useState(settings);
  const [version, setVersion] = useState(<></>);
  const { setError } = useGlobalErrorStore((x) => x);
  const navigate = useNavigate();

  useEffect(() => {
    getSystemClient()
      .getVersion({})
      .then((resp) => {
        setVersion(
          resp.response.version.split("\n").reduce(
            (a, b, _) => {
              return (
                <>
                  {a}
                  <br />
                  {b}
                </>
              );
            },
            <></>,
          ),
        );
      })
      .catch((error: unknown) => {
        setError({
          body: String(error),
        });
      });
  }, [setError]);

  useEffect(() => {
    setDirtySettings(settings);
  }, [settings]);

  return (
    <>
      <div className="flex flex-row gap-2 items-center">
        <Settings />
        <H2>Settings</H2>
      </div>

      <Checkbox
        checked={dirtySettings.darkMode}
        description="Whether to use dark mode (checked), or light mode (unchecked) for RPG Book."
        label="Dark Mode"
        onChange={(event) => {
          const s = structuredClone(dirtySettings);
          s.darkMode = event.currentTarget.checked;
          setDirtySettings(s);
        }}
      />

      <div className="flex flex-row justify-between gap-3">
        <Button
          onClick={() => {
            navigate("/");
          }}
          variant="default"
        >
          Cancel
        </Button>
        <Button
          onClick={() => {
            getSystemClient()
              .setSettings(dirtySettings)
              .then(() => {
                setSettings(dirtySettings);
                navigate("/");
              })
              .catch((error: unknown) => {
                setError({
                  body: String(error),
                });
              });
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
          checked={dirtySettings.devMode}
          description="Do not turn developer mode on unless you know what you are doing."
          label="Developer Mode"
          onChange={(event) => {
            const s = structuredClone(dirtySettings);
            s.devMode = event.currentTarget.checked;
            setDirtySettings(s);
          }}
        />
      </div>

      <InputDescription>{version}</InputDescription>
    </>
  );
}
