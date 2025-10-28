import { type ReactNode } from "react";
import { P } from "@/components/typography/P";

interface Props {
  no: ReactNode;
  onNo: () => void;
  onYes: () => void;
  yes: ReactNode;
}

// Anything with two non-mantine button yes no options, i.e: edit, delete
export function YesNo(props: Readonly<Props>): ReactNode {
  return (
    <div className="flex flex-row gap-3 self-start justify-end-safe">
      <button
        className="cursor-pointer h-fit"
        onClick={() => {
          props.onYes();
        }}
      >
        <P className="flex flex-row gap-1">{props.yes}</P>
      </button>
      <button
        className="cursor-pointer h-fit"
        onClick={() => {
          props.onNo();
        }}
      >
        <P className="flex flex-row gap-1 text-red-500 font-bold">{props.no}</P>
      </button>
    </div>
  );
}
