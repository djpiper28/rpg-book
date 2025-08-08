import { Button, TextInput } from "@mantine/core";
import { Plus, Text } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { CreateFileInput } from "@/components/input/createFileInput";
import { H2 } from "@/components/typography/H2";
import { H3 } from "@/components/typography/H3";
import { DbExtension } from "@/lib/databaseTypes";
import { projectClient } from "@/lib/grpcClient/client";
import { randomPlace } from "@/lib/random/randomProjectName";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useTabStore } from "@/stores/tabStore";
import { projectPath } from "../project/path";

function validateProjectName(rawName: string): string {
  const name = rawName.trim().normalize("NFC");

  if (name === "") {
    return "Please chose a name for your project";
  }

  return "";
}

export function Component() {
  const tabs = useTabStore((x) => x);
  const navigate = useNavigate();
  const [projectName, setProjectName] = useState<string>(randomPlace());
  const [saveLocation, setSaveLocation] = useState<string | null>();
  const { setError } = useGlobalErrorStore((x) => x);

  useEffect(() => {
    if (validateProjectName(projectName) !== "") {
      setSaveLocation(undefined);
      return;
    }

    const filteredName = projectName
      .trim()
      .normalize("NFC")
      .replace("/", "")
      .slice(0, 30);

    setSaveLocation(`${filteredName}${DbExtension}`);
  }, [projectName]);

  return (
    <>
      <H2>Create A New Project</H2>
      <TextInput
        description="A cool, user friendly name for your project."
        error={validateProjectName(projectName)}
        label="Project Name"
        onChange={(ev) => {
          const name = ev.target.value;
          setProjectName(name);
        }}
        placeholder="Project name"
        rightSection={<Text />}
        value={projectName}
        withAsterisk={true}
      />

      <H3>Advanced Settings</H3>
      <CreateFileInput
        description="Where to save the file, you usually don't need to touch this."
        error={saveLocation ? "" : "Please chose a location to save as"}
        filters={[
          {
            extensions: [DbExtension],
            name: "Project (*.sqlite)",
          },
        ]}
        label="Save Location"
        onChange={(f) => {
          if (!f.endsWith(DbExtension)) {
            f += DbExtension;
          }

          setSaveLocation(f);
        }}
        placeholder="Save location"
        // rightSection={<FileIcon />}
        value={saveLocation ?? ""}
        withAsterisk={true}
      />

      <Button
        onClick={() => {
          if (validateProjectName(projectName) !== "") {
            return;
          }

          if (!saveLocation) {
            return;
          }

          projectClient
            .createProject({
              fileName: saveLocation,
              projectName,
            })
            .then(async (resp) => {
              tabs.addTab({ id: resp.response.id }, projectName);
              await navigate(projectPath);
            })
            .catch((error: unknown) => {
              setError({
                body: String(error),
              });
            });
        }}
      >
        <Plus />
        Create Project
      </Button>
    </>
  );
}
