import { Pencil, Trash } from "lucide-react";
import { type ReactNode } from "react";
import { P } from "../typography/P";
import { YesNo } from "./yesNo";

interface Props {
  delete: () => void;
  edit: () => void;
}

export function EditDelete(props: Readonly<Props>): ReactNode {
  return (
    <YesNo
      no={
        <P className="flex flex-row gap-1 text-red-500 font-bold">
          <Trash />
          Delete
        </P>
      }
      onNo={() => {
        props.delete();
      }}
      onYes={() => {
        props.edit();
      }}
      yes={
        <P className="flex flex-row gap-1">
          <Pencil />
          Edit
        </P>
      }
    />
  );
}
