import Markdown from "react-markdown";

interface Props {
  markdown: string;
}

export default function MarkdownRenderer(props: Readonly<Props>) {
  return (
    <div className="border border-gray-500 rounded-lg p-2 min-w-50 min-h-10">
      <Markdown disallowedElements={["a"]} skipHtml={true}>
        {props.markdown}
      </Markdown>
    </div>
  );
}
