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
export const AgentConfigName = {
    AgentConfigNameIdentity: 'identity'
} as const;
export type AgentConfigName = typeof AgentConfigName[keyof typeof AgentConfigName];


export function AgentConfigNameFromJSON(json: any): AgentConfigName {
    return AgentConfigNameFromJSONTyped(json, false);
}

export function AgentConfigNameFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfigName {
    return json as AgentConfigName;
}

export function AgentConfigNameToJSON(value?: AgentConfigName | null): any {
    return value as any;
}

