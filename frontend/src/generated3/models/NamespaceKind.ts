/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.1
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
export const NamespaceKind = {
    NamespaceKindSystem: 'sys',
    NamespaceKindProfile: 'profile',
    NamespaceKindCaRoot: 'ca-root',
    NamespaceKindCaInt: 'ca-int',
    NamespaceKindServicePrincipal: 'service-principal',
    NamespaceKindGroup: 'group',
    NamespaceKindDevice: 'device',
    NamespaceKindUser: 'user',
    NamespaceKindApplication: 'application'
} as const;
export type NamespaceKind = typeof NamespaceKind[keyof typeof NamespaceKind];


export function NamespaceKindFromJSON(json: any): NamespaceKind {
    return NamespaceKindFromJSONTyped(json, false);
}

export function NamespaceKindFromJSONTyped(json: any, ignoreDiscriminator: boolean): NamespaceKind {
    return json as NamespaceKind;
}

export function NamespaceKindToJSON(value?: NamespaceKind | null): any {
    return value as any;
}

