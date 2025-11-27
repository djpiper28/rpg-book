import { Input } from "@mantine/core";
import { type ReactNode } from "react";
import { P } from "@/components/typography/P";

interface Props<T> {
  error?: string;
  onChange: (text: string) => void;
  placeholder: string;
  render: (element: T) => ReactNode;
  rightElement?: ReactNode;
  searchRes: T[];
}

export function Search<T>(props: Readonly<Props<T>>): ReactNode {
  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-row gap-2 justify-between">
        <Input
          className="flex-grow"
          onChange={(event) => {
            props.onChange(event.target.value);
          }}
          placeholder={props.placeholder}
        />
        {props.rightElement}
      </div>
      {props.error && (
        <P className="text-red-500">Search error: {props.error}</P>
      )}
      {props.searchRes.map((x) => {
        return props.render(x);
      })}
    </div>
  );
}
