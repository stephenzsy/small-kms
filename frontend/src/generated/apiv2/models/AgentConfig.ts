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

import {
    AgentConfigIdentity,
    instanceOfAgentConfigIdentity,
    AgentConfigIdentityFromJSON,
    AgentConfigIdentityFromJSONTyped,
    AgentConfigIdentityToJSON,
} from './AgentConfigIdentity';

/**
 * @type AgentConfig
 * 
 * @export
 */
export type AgentConfig = { name: 'identity' } & AgentConfigIdentity;

export function AgentConfigFromJSON(json: any): AgentConfig {
    return AgentConfigFromJSONTyped(json, false);
}

export function AgentConfigFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfig {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    switch (json['name']) {
        case 'identity':
            return {...AgentConfigIdentityFromJSONTyped(json, true), name: 'identity'};
        default:
            throw new Error(`No variant of AgentConfig exists with 'name=${json['name']}'`);
    }
}

export function AgentConfigToJSON(value?: AgentConfig | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    switch (value['name']) {
        case 'identity':
            return AgentConfigIdentityToJSON(value);
        default:
            throw new Error(`No variant of AgentConfig exists with 'name=${value['name']}'`);
    }

}
