import { useBoolean } from "ahooks";
import { useContext } from "react";
import { useParams } from "react-router-dom";
import AlertDialog from "../components/AlertDialog";
import { Button } from "../components/Button";
import { PolicyContext, PolicyContextProvider } from "./PolicyContext";

function CurrentPolicy() {
  const { policy } = useContext(PolicyContext);
  return (
    <div className="bg-white p-6">
      <h1 className="text-2xl font-semibold">Current policy</h1>
      {policy === undefined ? (
        <div>Loading...</div>
      ) : policy === null ? (
        <div>No policy found</div>
      ) : (
        <pre>{JSON.stringify(policy, undefined, 2)}</pre>
      )}
    </div>
  );
}

function CurrentPolicyControl() {
  const { policy, deletePolicy, applyPolicy } = useContext(PolicyContext);
  const [
    disableDialogOpen,
    { setFalse: closeDisableDialog, setTrue: openDisableDialog },
  ] = useBoolean(false);
  const [
    deleteDialogOpen,
    { setFalse: closeDeleteDialog, setTrue: openDeleteDialog },
  ] = useBoolean(false);
  if (!policy) return null;
  if (policy.deleted)
    return (
      <>
        <Button color="danger" variant="soft" onClick={openDeleteDialog}>
          Delete
        </Button>
        <AlertDialog
          title="Delete policy"
          description="Are you sure you want to delete this policy? This action cannot be undone."
          open={deleteDialogOpen}
          onClose={closeDeleteDialog}
          okText="Delete policy"
          onOk={() => {
            deletePolicy(true);
            closeDeleteDialog();
          }}
        />
      </>
    );
  return (
    <>
      <Button color="danger" variant="soft" onClick={openDisableDialog}>
        Disable
      </Button>
      <AlertDialog
        title="Disable policy"
        description="Are you sure you want to disable this policy?"
        open={disableDialogOpen}
        onClose={closeDisableDialog}
        okText="Disable policy"
        onOk={() => {
          deletePolicy();
          closeDisableDialog();
        }}
      />
      <Button variant="primary" onClick={applyPolicy}>
        Apply policy
      </Button>
    </>
  );
}

export default function PolicyPage() {
  const { namespaceId, policyId } = useParams();

  return (
    <>
      <h1 className="text-4xl font-semibold">Policy</h1>
      <PolicyContextProvider namespaceId={namespaceId!} policyId={policyId!}>
        <CurrentPolicy />
        <div className="flex flex-row items-center gap-x-6">
          <CurrentPolicyControl />
        </div>
      </PolicyContextProvider>
    </>
  );
}
