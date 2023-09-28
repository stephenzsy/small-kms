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
export const JwkKeySize = {
    KeySize_2048: 2048,
    KeySize_3072: 3072,
    KeySize_4096: 4096
} as const;
export type JwkKeySize = typeof JwkKeySize[keyof typeof JwkKeySize];


export function JwkKeySizeFromJSON(json: any): JwkKeySize {
    return JwkKeySizeFromJSONTyped(json, false);
}

export function JwkKeySizeFromJSONTyped(json: any, ignoreDiscriminator: boolean): JwkKeySize {
    return json as JwkKeySize;
}

export function JwkKeySizeToJSON(value?: JwkKeySize | null): any {
    return value as any;
}
