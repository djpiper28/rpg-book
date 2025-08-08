import { Button, Textarea } from "@mantine/core";
import { useState } from "react";
import MarkdownRenderer from "../renderers/markdown";
import { P } from "../typography/P";

interface Props {
  label: string;
  setValue: (value: string) => void;
  value: string;
}

export function MarkdownEditor(props: Readonly<Props>) {
  const [showPreview, setShowPreview] = useState(false);

  return (
    <div className="flex flex-row gap-5 justify-between">
      <div className="flex flex-col gap-1 flex-grow">
        <Textarea
          autosize
          className="grow min-w-50"
          label={props.label}
          minRows={10}
          onChange={(event) => {
            props.setValue(event.target.value);
          }}
          placeholder={`# Character description
---
Supports **markdown**.`}
          resize="vertical"
          value={props.value}
        />
        <Button
          className="grow-0"
          onClick={() => {
            setShowPreview(!showPreview);
          }}
          variant="subtle"
        >
          {showPreview ? "Hide" : "Show"} preview
        </Button>
      </div>
      {showPreview && (
        <div className="flex flex-col flex-grow">
          <P>Preview</P>
          <MarkdownRenderer markdown={props.value} />
        </div>
      )}
    </div>
  );
}
