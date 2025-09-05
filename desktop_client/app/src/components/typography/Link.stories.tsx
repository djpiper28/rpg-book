import { type Meta, type StoryObj } from "@storybook/react-vite";
import { Link } from "./Link";

const meta = {
  component: Link,
} satisfies Meta<typeof Link>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    children: "Lorem ipsum",
    href: "https://lichess.org",
  },
};
