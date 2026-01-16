import { Table } from "@mantine/core";
import { type ReactNode } from "react";
import { H1 } from "@/components/typography/H1";
import { H2 } from "@/components/typography/H2";
import { H3 } from "@/components/typography/H3";
import { Link } from "@/components/typography/Link";
import { P } from "@/components/typography/P";

export function Component(): ReactNode {
  return (
    <>
      <H1>Search Help</H1>
      <H2>General Syntax</H2>
      <P>
        RPG Book uses syntax similar to Scryfall, and Cockatrice (Magic the
        Gathering apps). It is easy to use and fairly powerful. The most simple
        way to search is to put parts of the name of what you wish to search
        into the search bar, i.e: &quot;David&quot;.
      </P>
      <P>
        More advanced queries are possible, the syntax lets you chose a field to
        search (i.e: name, description, etc... - see{" "}
        <a href="#scopes">scopes</a> for details). You can then use a series of
        operators to search the determine how the field should be searched.
      </P>
      <P>
        Typically a query will something look like this: &quot;name:dave and
        (desc:dodgy or desc:&quot;brexit vote&quot;)&quot;. Each term follows
        the syntax of &quot;&lt;field&tg;&lt;
        <a href="set-generation-operators">operator</a>&gt;&lt;value&gt;&quot;.
      </P>
      <P>
        The different terms can be combined using &quot;and&quot;, &amp;
        &quot;or&quot; and brackets can be used to remove ambiguity.
      </P>

      <H2 id="set-generation-operators">Set Generation Operators</H2>
      <P>
        All of the string operators are case insensitive and accent permissive.
      </P>
      <Table variant="vertical">
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Operator</Table.Th>
            <Table.Th>Name</Table.Th>
            <Table.Th>Usage</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tr>
          <Table.Th>:</Table.Th>
          <Table.Th>Includes Operator</Table.Th>
          <Table.Th>
            Checks if a string (the value) is in the field that you are
            searching. For example &quot;name:dav&quot; will match names such as
            Dave, Davros, David, McDavidson, etc...
          </Table.Th>
        </Table.Tr>
        <Table.Tr>
          <Table.Th>=</Table.Th>
          <Table.Th>Equals Operator</Table.Th>
          <Table.Th>
            Checks if a string (the value) is equal to the field that you are
            searching. For example &quot;name=&quot;David Cameron&quot;&quot;
            will only ever match David Cameron, like all operators it is case
            insensitive and has accent permissive.
            <br />
            This is also valid for numerical searches where it checks that the
            numbers are equal.
          </Table.Th>
        </Table.Tr>

        <Table.Tr>
          <Table.Th>&lt;</Table.Th>
          <Table.Th>Less Than Operator</Table.Th>
          <Table.Th>
            Checks if a numerical field is less than the value.
          </Table.Th>
        </Table.Tr>
        <Table.Tr>
          <Table.Th>&lt;=</Table.Th>
          <Table.Th>Less Than Or Equal Operator</Table.Th>
          <Table.Th>
            Checks if a numerical field is less than of equal to the value.
          </Table.Th>
        </Table.Tr>
        <Table.Tr>
          <Table.Th>&gt;</Table.Th>
          <Table.Th>Greater Than Operator</Table.Th>
          <Table.Th>
            Checks if a numerical field is greater than the value.
          </Table.Th>
        </Table.Tr>
        <Table.Tr>
          <Table.Th>&gt;=</Table.Th>
          <Table.Th>Greater Than Or Equal Operator</Table.Th>
          <Table.Th>
            Checks if a numerical field is greater than of equal to the value.
          </Table.Th>
        </Table.Tr>
      </Table>

      <H2 id="scopes">Search Scopes</H2>
      <P>
        There are a series of different scopes for searching, think of these as
        each of the different search bars in the RPG Book.
      </P>
      <H3 id="#characters">Characters</H3>
      <Table variant="vertical">
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Field</Table.Th>
            <Table.Th>Description</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tr>
          <Table.Th>name</Table.Th>
          <Table.Th>The name of the character i.e: George W Bush</Table.Th>
        </Table.Tr>
        <Table.Tr>
          <Table.Th>desc, or description</Table.Th>
          <Table.Th>
            The description of the character i.e: Likes to eat mushy peas on a
            Wednesday with his wife.
          </Table.Th>
        </Table.Tr>
      </Table>

      <H3>Notes</H3>
      <Table variant="vertical">
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Field</Table.Th>
            <Table.Th>Description</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tr>
          <Table.Th>name</Table.Th>
          <Table.Th>The name of the note i.e: Map of Scary Dungeon</Table.Th>
        </Table.Tr>
        <Table.Tr>
          <Table.Th>markdown, contents, desc, or description</Table.Th>
          <Table.Th>
            The text contents of description of the note i.e: 2 spuds
          </Table.Th>
        </Table.Tr>
        <Table.Tr>
          <Table.Th>character.*</Table.Th>
          <Table.Th>
            Allows you to search by the fields of related characters i.e:
            character.name=&quot;Dave&qauot;. See{" "}
            <Link href="#characters">characters</Link>.
          </Table.Th>
        </Table.Tr>
      </Table>
    </>
  );
}
