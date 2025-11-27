import { Input } from "@mantine/core";
import { HelpCircleIcon } from "lucide-react";
import { type ReactNode } from "react";
import { P } from "@/components/typography/P";
import { searchHelpPath } from "@/pages/help/search/path";
import { useNavigate } from "react-router";

interface Props<T> {
  elementWrapper: (children: ReactNode[]) => ReactNode;
  error?: string;
  onChange: (text: string) => void;
  placeholder: string;
  render: (element: T) => ReactNode;
  rightElement?: ReactNode;
  searchRes: T[];
}

export function Search<T>(props: Readonly<Props<T>>): ReactNode {
  const navigate = useNavigate();

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
        <button
          className="cursor-pointer"
          onClick={() => {
            navigate(searchHelpPath);
          }}
        >
          <HelpCircleIcon />
        </button>
        {props.rightElement}
      </div>
      {props.error && (
        <P className="text-red-500">Search error: {props.error}</P>
      )}
      {props.elementWrapper(
        props.searchRes.map((x) => {
          return props.render(x);
        }),
      )}
    </div>
  );
}
