import { type Meta, type StoryObj } from "@storybook/react-vite";
import { MarkdownEditor } from "./markdownEditor";

const meta = {
  component: MarkdownEditor,
} satisfies Meta<typeof MarkdownEditor>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    label: "Test",
    setValue: (value: string) => {
      console.log(value);
    },
    value: `# Testing 123
## This is a test
Hello world !!!

---

[uwu](http://google.com/?q=uwu)"
`,
  },
};
