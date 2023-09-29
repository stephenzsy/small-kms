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
export const JwkKeyOperation = {
    KeyOpSign: 'sign',
    KeyOpVerify: 'verify',
    KeyOpEncrypt: 'encrypt',
    KeyOpDecrypt: 'decrypt',
    KeyOpWrapKey: 'wrapKey',
    KeyOpUnwrapKey: 'unwrapKey'
} as const;
export type JwkKeyOperation = typeof JwkKeyOperation[keyof typeof JwkKeyOperation];


export function JwkKeyOperationFromJSON(json: any): JwkKeyOperation {
    return JwkKeyOperationFromJSONTyped(json, false);
}

export function JwkKeyOperationFromJSONTyped(json: any, ignoreDiscriminator: boolean): JwkKeyOperation {
    return json as JwkKeyOperation;
}

export function JwkKeyOperationToJSON(value?: JwkKeyOperation | null): any {
    return value as any;
}

