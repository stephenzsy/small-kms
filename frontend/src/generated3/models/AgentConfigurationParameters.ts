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

import {
    AgentConfigurationAgentActiveHostBootstrap,
    instanceOfAgentConfigurationAgentActiveHostBootstrap,
    AgentConfigurationAgentActiveHostBootstrapFromJSON,
    AgentConfigurationAgentActiveHostBootstrapFromJSONTyped,
    AgentConfigurationAgentActiveHostBootstrapToJSON,
} from './AgentConfigurationAgentActiveHostBootstrap';
import {
    AgentConfigurationAgentActiveServer,
    instanceOfAgentConfigurationAgentActiveServer,
    AgentConfigurationAgentActiveServerFromJSON,
    AgentConfigurationAgentActiveServerFromJSONTyped,
    AgentConfigurationAgentActiveServerToJSON,
} from './AgentConfigurationAgentActiveServer';

/**
 * @type AgentConfigurationParameters
 * 
 * @export
 */
export type AgentConfigurationParameters = { name: 'agent-active-host-bootstrap' } & AgentConfigurationAgentActiveHostBootstrap | { name: 'agent-active-server' } & AgentConfigurationAgentActiveServer;

export function AgentConfigurationParametersFromJSON(json: any): AgentConfigurationParameters {
    return AgentConfigurationParametersFromJSONTyped(json, false);
}

export function AgentConfigurationParametersFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfigurationParameters {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    switch (json['name']) {
        case 'agent-active-host-bootstrap':
            return {...AgentConfigurationAgentActiveHostBootstrapFromJSONTyped(json, true), name: 'agent-active-host-bootstrap'};
        case 'agent-active-server':
            return {...AgentConfigurationAgentActiveServerFromJSONTyped(json, true), name: 'agent-active-server'};
        default:
            throw new Error(`No variant of AgentConfigurationParameters exists with 'name=${json['name']}'`);
    }
}

export function AgentConfigurationParametersToJSON(value?: AgentConfigurationParameters | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    switch (value['name']) {
        case 'agent-active-host-bootstrap':
            return AgentConfigurationAgentActiveHostBootstrapToJSON(value);
        case 'agent-active-server':
            return AgentConfigurationAgentActiveServerToJSON(value);
        default:
            throw new Error(`No variant of AgentConfigurationParameters exists with 'name=${value['name']}'`);
    }

}
