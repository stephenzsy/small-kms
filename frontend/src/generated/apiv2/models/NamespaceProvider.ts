/* tslint:disable */
/* eslint-disable */
/**
 * Cryptocat API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.3
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
export const NamespaceProvider = {
    NamespaceProviderProfile: 'profile',
    NamespaceProviderAgent: 'agent',
    NamespaceProviderServicePrincipal: 'service-principal'
} as const;
export type NamespaceProvider = typeof NamespaceProvider[keyof typeof NamespaceProvider];


export function NamespaceProviderFromJSON(json: any): NamespaceProvider {
    return NamespaceProviderFromJSONTyped(json, false);
}

export function NamespaceProviderFromJSONTyped(json: any, ignoreDiscriminator: boolean): NamespaceProvider {
    return json as NamespaceProvider;
}

export function NamespaceProviderToJSON(value?: NamespaceProvider | null): any {
    return value as any;
}

