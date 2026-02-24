import { type ReactNode } from "react";
import Markdown from "react-markdown";
import { H1 } from "../typography/H1";
import { H2 } from "../typography/H2";
import { H3 } from "../typography/H3";
import { Link } from "../typography/Link";
import { P } from "../typography/P";

interface SafeLinkProps {
  children?: string | ReactNode;
  href?: string;
}

function SafeLink(props: Readonly<SafeLinkProps>): ReactNode {
  return (
    <Link href={props.href ?? ""} openInBrowser={true}>
      {props.children ?? ""}
    </Link>
  );
}

interface Props {
  className?: string;
  markdown: string;
}

export default function MarkdownRenderer(props: Readonly<Props>): ReactNode {
  return (
    <div
      className={`border border-gray-500 rounded-lg p-2 min-w-50 min-h-10 w-full h-full ${props.className ?? ""}`}
    >
      <Markdown
        components={{
          a: SafeLink,
          h1: H1,
          h2: H2,
          h3: H3,
          p: P,
        }}
        skipHtml={true}
      >
        {props.markdown}
      </Markdown>
    </div>
  );
}
