import { Button, Table } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import { useNavigate, useRouteError } from "react-router";
import { H1 } from "./components/typography/H1";
import { H2 } from "./components/typography/H2";
import { P } from "./components/typography/P";
import { getBuildType, getSystemVersion } from "./lib/electron/index.ts";
import { getLogger, getSystemClient } from "./lib/grpcClient/client";
import { mustVoid } from "./lib/utils/errorHandlers";
import { indexPath } from "./pages/path";

const commonCss = "text-wrap text-red-400";

export function ErrorPage(): ReactNode {
  const navigate = useNavigate();
  const routeError = useRouteError();
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  const routeErrCast = routeError as Error;
  const [version, setVersion] = useState(<></>);
  const [backendError, setBackendError] = useState("");

  useEffect(() => {
    getSystemClient()
      .getVersion({})
      .then((resp) => {
        setVersion(
          resp.response.version.split("\n").reduce(
            (a, b, _) => {
              return (
                <>
                  {a}
                  <br />
                  {b}
                </>
              );
            },
            <></>,
          ),
        );
      })
      .catch((error_: unknown) => {
        // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
        setBackendError(`${error_}`);
      });
  }, []);

  getLogger().error("Error boundary triggered", {
    backendError,
    location: globalThis.location.href,
    // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
    routeError: `${routeError}`,
  });

  return (
    <div className="flex flex-col gap-3 p-10">
      <H1>An Error Has Occured</H1>
      <P>You may need to restart the application</P>
      <Button
        className="w-min"
        onClick={() => {
          mustVoid(navigate(indexPath));
        }}
      >
        Try returning to home page
      </Button>

      <H2>Debug Information For Github Issues</H2>
      <Table variant="vertical">
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Key</Table.Th>
            <Table.Th>Value</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          <Table.Tr>
            <Table.Th>Current page</Table.Th>
            <Table.Th>
              <pre className={commonCss}>{globalThis.location.href}</pre>
            </Table.Th>
          </Table.Tr>
          <Table.Tr>
            <Table.Th>Route Error</Table.Th>
            <Table.Th>
              <pre className={commonCss}>
                {JSON.stringify(routeErrCast.message)}
              </pre>
            </Table.Th>
          </Table.Tr>
          <Table.Tr>
            <Table.Th>Backend Call Error</Table.Th>
            <Table.Th>
              <pre className={commonCss}>{backendError}</pre>
            </Table.Th>
          </Table.Tr>
          <Table.Tr>
            <Table.Th>Version</Table.Th>
            <Table.Th>
              <pre className={commonCss}>{version}</pre>
            </Table.Th>
          </Table.Tr>
          <Table.Tr>
            <Table.Th>OS Version</Table.Th>
            <Table.Th>
              <pre className={commonCss}>{getSystemVersion()}</pre>
            </Table.Th>
          </Table.Tr>
          <Table.Tr>
            <Table.Th>Build Type</Table.Th>
            <Table.Th>
              <pre className={commonCss}>{getBuildType()}</pre>
            </Table.Th>
          </Table.Tr>
        </Table.Tbody>
      </Table>
    </div>
  );
}
