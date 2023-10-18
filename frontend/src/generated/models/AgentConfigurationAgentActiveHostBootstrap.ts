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

import { exists, mapValues } from '../runtime';
import type { AgentConfigName } from './AgentConfigName';
import {
    AgentConfigNameFromJSON,
    AgentConfigNameFromJSONTyped,
    AgentConfigNameToJSON,
} from './AgentConfigName';
import type { AgentConfigurationActiveHostControllerContainer } from './AgentConfigurationActiveHostControllerContainer';
import {
    AgentConfigurationActiveHostControllerContainerFromJSON,
    AgentConfigurationActiveHostControllerContainerFromJSONTyped,
    AgentConfigurationActiveHostControllerContainerToJSON,
} from './AgentConfigurationActiveHostControllerContainer';

/**
 * 
 * @export
 * @interface AgentConfigurationAgentActiveHostBootstrap
 */
export interface AgentConfigurationAgentActiveHostBootstrap {
    /**
     * 
     * @type {AgentConfigName}
     * @memberof AgentConfigurationAgentActiveHostBootstrap
     */
    name: AgentConfigName;
    /**
     * 
     * @type {AgentConfigurationActiveHostControllerContainer}
     * @memberof AgentConfigurationAgentActiveHostBootstrap
     */
    controllerContainer: AgentConfigurationActiveHostControllerContainer;
}

/**
 * Check if a given object implements the AgentConfigurationAgentActiveHostBootstrap interface.
 */
export function instanceOfAgentConfigurationAgentActiveHostBootstrap(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "controllerContainer" in value;

    return isInstance;
}

export function AgentConfigurationAgentActiveHostBootstrapFromJSON(json: any): AgentConfigurationAgentActiveHostBootstrap {
    return AgentConfigurationAgentActiveHostBootstrapFromJSONTyped(json, false);
}

export function AgentConfigurationAgentActiveHostBootstrapFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfigurationAgentActiveHostBootstrap {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'name': AgentConfigNameFromJSON(json['name']),
        'controllerContainer': AgentConfigurationActiveHostControllerContainerFromJSON(json['controllerContainer']),
    };
}

export function AgentConfigurationAgentActiveHostBootstrapToJSON(value?: AgentConfigurationAgentActiveHostBootstrap | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'name': AgentConfigNameToJSON(value.name),
        'controllerContainer': AgentConfigurationActiveHostControllerContainerToJSON(value.controllerContainer),
    };
}
