import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
} from "@radix-ui/react-navigation-menu";
import { H1 } from "./components/typography/H1";

function App() {
  return (
    <div className="flex flex-col gap-3">
      <NavigationMenu>
        <NavigationMenuList>
          <NavigationMenuItem>
            <H1>RPG Book</H1>
            <NavigationMenuTrigger>Projects</NavigationMenuTrigger>
            <NavigationMenuContent>
              <NavigationMenuLink>Create a project</NavigationMenuLink>
            </NavigationMenuContent>
            <NavigationMenuTrigger>Settings</NavigationMenuTrigger>
            <NavigationMenuContent>
              <NavigationMenuLink>Appearance</NavigationMenuLink>
            </NavigationMenuContent>
          </NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>

      <p className="text-red-500">TEST</p>
    </div>
  );
}

export default App;
