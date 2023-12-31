import { NIL } from "uuid";

export const enum WellknownId {
  nsRootCa = "00000000-0000-0000-0000-000000000001",
  nsTestRootCa = "00000000-0000-0000-0000-00000000000f",

  nsIntCaIntranet = "00000000-0000-0000-0000-000000000012",
  nsTestIntCa = "00000000-0000-0000-0000-0000000000ff",

  defaultPolicyIdCertRequest = "00000000-0000-0000-0001-000000000001",
  defaultPolicyIdCertEnroll = "00000000-0000-0000-0001-000000000002",
  defaultPolicyIdAadAppCred = "00000000-0000-0000-0001-000000000003",
}

export const uuidNil = NIL;
