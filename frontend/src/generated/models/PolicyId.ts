/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


/**
 * 
 * @export
 */
export const PolicyId = {
    PolicyID_CertEnroll: 'cert-enrollment'
} as const;
export type PolicyId = typeof PolicyId[keyof typeof PolicyId];


export function PolicyIdFromJSON(json: any): PolicyId {
    return PolicyIdFromJSONTyped(json, false);
}

export function PolicyIdFromJSONTyped(json: any, ignoreDiscriminator: boolean): PolicyId {
    return json as PolicyId;
}

export function PolicyIdToJSON(value?: PolicyId | null): any {
    return value as any;
}

