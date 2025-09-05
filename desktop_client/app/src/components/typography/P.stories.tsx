import { type Meta, type StoryObj } from "@storybook/react-vite";
import { P } from "./P";

const meta = {
  component: P,
} satisfies Meta<typeof P>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    children: "Lorem ipsum",
  },
};
