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
export const NamespaceType = {
    NamespaceType_BuiltIn: 'builtIn',
    NamespaceType_MsGraphServicePrincipal: '#microsoft.graph.servicePrincipal',
    NamespaceType_MsGraphUser: '#microsoft.graph.user'
} as const;
export type NamespaceType = typeof NamespaceType[keyof typeof NamespaceType];


export function NamespaceTypeFromJSON(json: any): NamespaceType {
    return NamespaceTypeFromJSONTyped(json, false);
}

export function NamespaceTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): NamespaceType {
    return json as NamespaceType;
}

export function NamespaceTypeToJSON(value?: NamespaceType | null): any {
    return value as any;
}

