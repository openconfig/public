

## Controller Card Power Control

This document describes operational use cases, rules and telemetry for power control of `CONTROLLER-CARD` components in network devices using OpenConfig.


## Operational use cases for CONTROLLER_CARD power control

1. Operator suspects a network device issue is occurring due to faults on a CONTROLLER_CARD. The operator wants to quickly disable the CONTROLLER_CARD and later, troubleshoot the issue in a maintenance window. The operator also wants to ensure that the card stays disabled even if the device is rebooted. Some example scenarios below:
   
   a. It is possible that the faulty controller card is in a boot loop and the operator finds it best to disable the card to prevent undesired state. Expectations are also that the power remains disabled  even after a system reboot.
   
   b. CPU overheating, memory errors or other hardware problems with the PRIMARY card may require the operator to proactively power-off the card and not let it be online during regular operation until replacement of the hardware is completed.
   
2. Operator feels it important to keep the card shutdown to prevent unexpected outcomes post the physical swap of the faulty/alerting card so the operator can online the card in a controlled environment post replacement.


## CONTROLLER_CARD power and redundancy requirements 

Electrical power for a CONTROLLER_CARD can be configured off using the OC-Path [/components/component/controller-card/config/power-admin-state](https://openconfig.net/projects/models/schemadocs/yangdoc/openconfig-platform.html#components-component-controller-card-config-power-admin-state) and setting its value to POWER\_DISABLED.  The following rules regarding CONTROLLER_CARD redundancy and changes to `power-admin-state` must be followed by the device:



1. Only a CONTROLLER_CARD in `state/redundant-role` 'SECONDARY' will honor a change in `config/power-admin-state` to POWER_DISABLED. If the controller-card's `redundant-role` is 'PRIMARY', and its `config/power-admin-state` is set to 'POWER_DISABLED', the NOS must allow the configuration. However, the `state/power-admin-state` must remain as 'POWER_ENABLED'.  A change in `state/power-admin-state` must take effect only on the next reboot or if the CONTROLLER_CARD is 'SECONDARY'. Examples of scenarios include:
             
### Scenario 1 - Power off a secondary CONTROLLER_CARD

Let's say controller-card0 is PRIMARY and controller-card1 is SECONDARY and controller-card1 receives an operator driven change of config/power-admin-state = POWER_DISABLED, then controller-card0 will power-off controller-card1 immediately.  The leaf state/power-admin-state for controller-card1 must also be POWER_DISABLED. The 'state/last-poweroff-reason/trigger' should show as USER_INITIATED. The NOS may optionally update state/last-poweroff-reason/details. state/last-poweroff-time should record the time when the card was powered-off. For example:

```
/components/component[name=controller0]/state/redundant-role, PRIMARY
/components/component[name=controller0]/controller-card/config/power-admin-state, POWER_ENABLED
/components/component[name=controller0]/controller-card/state/power-admin-state, POWER_ENABLED
/components/component[name=controller1]/state/redundant-role, SECONDARY
/components/component[name=controller1]/controller-card/config/power-admin-state, POWER_DISABLED
/components/component[name=controller1]/controller-card/state/power-admin-state, POWER_DISABLED
/components/component[name=controller1]/state/last-poweroff-reason/trigger, USER_INITIATED
/components/component[name=controller1]/state/last-poweroff-reason/details, "User initiated Shutdown"
/components/component[name=controller1]/state/last-poweroff-time, 1706590714681765937
```

### Scenario 2 - Power off primary CONTROLLER_CARD
    
1. If controller-card0 is PRIMARY and controller-card1 is SECONDARY and if controller-card0 is set to config/power-admin-state = POWER_DISABLED by an operator, then controller-card0 will stay powered-on until the next reboot. `state/power-admin-state` must show as POWER_ENABLED. If a reboot of the PRIMARY CONTROLLER_CARD occurs,  `state/last-poweroff-time` must record the time when the card powers-off, `state/last-poweroff-reason/trigger` must show as USER_INITIATED and `/state/last-poweroff-reason/details` may be updated. For example: 

When controller-card0 is PRIMARY:

```          
/components/component[name=controller0]/state/redundant-role, PRIMARY
/components/component[name=controller0]/controller-card/config/power-admin-state, POWER_DISABLED
/components/component[name=controller0]/controller-card/state/power-admin-state, POWER_ENABLED
/components/component[name=controller1]/state/redundant-role, SECONDARY
/components/component[name=controller1]/controller-card/state/power-admin-state, POWER_ENABLED
```

After controller-card0 transitions to redundant-role SECONDARY:

```
/components/component[name=controller0]/state/redundant-role, SECONDARY
/components/component[name=controller0]/controller-card/config/power-admin-state, POWER_DISABLED
/components/component[name=controller0]/controller-card/state/power-admin-state, POWER_DISABLED
/components/component[name=controller0]/state/last-poweroff-reason/trigger, USER_INITIATED
/components/component[name=controller0]/state/last-poweroff-reason/details, "User initiated Shutdown"
/components/component[name=controller0]/state/last-poweroff-time, 1706590714681765937
/components/component[name=controller1]/state/redundant-role, PRIMARY
/components/component[name=controller1]/controller-card/config/power-admin-state, POWER_ENABLED
/components/component[name=controller1]/controller-card/state/power-admin-state, POWER_ENABLED
```

2. A controller-card which is in state/redundant-role=SECONDARY and is config/power-admin-state=POWER_DISABLED must remain powered off, even after a reboot. It is possible that after a reboot, both the controller cards are powered-on. However as soon as the configuration is loaded, the system must power-off the subject controller-card.

3. A CONTROLLER_CARD in redundant-role SECONDARY and state/power-admin-state = POWER_DISABLED cannot transition to redundant-role PRIMARY. If the PRIMARY CONTROLLER_CARD goes down, the device will be offline.

4. When a controller card boots up and loads it’s configuration, power-admin-state may be set to POWER_DISABLED. The CONTROLLER_CARD must then power off and never enter into any controller card primary/secondary election process. This also means that an implementation shouldn't start a controller card election process until the configuration is loaded and consumed.

	**Note:-** If an implementation's architecture do not allow for controlling the order in which the configuration is loaded and the PRIMARY/SECONDARY election process, then this rule can be relaxed as long as the implementation has proper arrangements to power-off the controller-card with `config/power-admin-state=POWER_DISABLED` configuration during reboots (Warm/Cold). 
   

5. On boot (cold or warm), if the chassis has a single controller card and it is configured for config/power-admin-state=POWER_DISABLED, it Must continue with the boot process ignoring the configuration. As per Rule#1 above, the controller-card  Must have it's state/power-admin-state as POWER_ENABLED given that a single CONTROLLER_CARD will have it's state/redundant-role=PRIMARY. In this case too, the system must log a message using severity "Warning", to inform the Operator about the situation.

6. In a Dual controller-card scenario, if a config is pushed for config/power-admin-state=POWER_DISABLED for either both controller-cards simultaneously or for one of the controller-cards while the other controller-card in the system is already configured for config/power-admin-state=POWER_DISABLED, then the implementation Must fail the configuration commit operation with an error similar to: "Not allowed to have both controller-cards configured for power-admin-state = POWER_DISABLED"

## Flowchart on the Rules above:

![Overview of the expected behavior](https://github.com/openconfig/public/tree/master/doc/img/controller_card.png?raw=true)

## Concerns and possible failure scenarios

1. If a PRIMARY card malfunctions and ends up in a bootloop, would this approach help?

    **Response:**


    In this scenario the expectations are that the implementation takes steps to initiate a controller-card switchover operation. Therefore, the standby controller card takes over the PRIMARY role and the system stabilizes allowing for gRPC connections to be established. In this situation if the operator pushes a configuration to shutdown the SECONDARY card, the PRIMARY card must be able to power-off the SECONDARY controller-card. In this scenario, the implementation (depending on their architecture) may also initiate a shutdown of the faulty card from the new PRIMARY card.

2. Since the operation relies on configuration, it is possible that the failure scenario may kick in before the configuration takes effect post a reboot.

    **Response:**


	Response here is the same as "1" above.



3. Both controller-cards are functional and the the secondary controller card is attempted to be powered off.

    **Response:**


    The shutdown can be initiated by pushing the command [/components/component[controller-card#]/controller-card/config/power-admin-state = POWER\_DISABLED](https://openconfig.net/projects/models/schemadocs/yangdoc/openconfig-platform.html#components-component-controller-card-config-power-admin-state) to the box. As a result, the secondary controller card whose [/components/component/state/redundant-role](https://openconfig.net/projects/models/schemadocs/yangdoc/openconfig-platform.html#components-component-state-redundant-role) is SECONDARY will be shutdown and its [/components/component[name="my_secondary_controller-card"]/controller-card/state/power-admin-state](https://openconfig.net/projects/models/schemadocs/yangdoc/openconfig-platform.html#components-component-controller-card-state-power-admin-state) will change to "POWER\_DISABLED".  If the configuration is saved and when the device reboots, the subject controller-card is expected to remain powered off.

4. Both controller-cards are functional and the the primary controller card is attempted to be shutdown

    **Response:**


    Follow Rule#1 above

5. Only one controller-card is present and the card is attempted to be shutdown.

    **Response:**


    Since the lone controller-card is PRIMARY, please follow Rule#1 above. After reboot, please follow Rule#5

6. Faulty card is shutdown and then replaced with a working card but device still has the config of `/components/component[controller-card#]/controller-card/config/power-admin-state = POWER_DISABLED`

    **Response:**


    The new card must remain powered off.  Power should only be enabled on the new card by a configuration change to set `/components/component[controller-card#]/controller-card/config/power-admin-state = POWER_ENABLED`


7. Both controller-cards present during the initial boot operation. However, the controller card that is configured as "POWER\_DISABLED" boots up sooner than the other card and takes over the ACTIVE role.

    **Response:**


    Follow Rule#4 above. The only exception is when it is the only controller card in the chassis. In the latter case, Rule#5 above must be followed.


8. Say we have a situation whereby, A dual controller card device is functional, secondary is powered off by changing power-admin-state to POWER\_DISABLED. Later, the PRIMARY card is removed and the system goes offline. If the device is rebooted (for example via externally removing and reapplying power ), what should happen?

    **Response:**
   
    If during reboot the device continues to have a single controller-card, then Rule#5 above should be followed.

9. According to Rule#5, a controller-card must ignore the configuration of `config/power-admin-state=POWER_DISABLED` and must continue with its boot process if it is the only controller-card in the chassis. This also means that, under split-brain situations when the communication between the controller-cards is broken, the option to use `config/power-admin-state=POWER_DISABLED` for one of the cards may not help post reboot?

   **Response**

    Lets walk through this situation to set the right context. We have a dual controller-card system which due to some reason (Hardware/Soaftware failure) got in to a split-brain scenaio. As per the question here, the operator decides to use `config/power-admin-state=POWER_DISABLED` command on one of the cards, somehow intiatiates a reboot process (Cold or Warm) and expects that the controller-card with `config/power-admin-state=POWER_DISABLED` stays disabled post reboot. There are several assumptions and nuances to this situation:

   a. One of the assumptions is that, the device is in a split-brain situation and would still allow connections for a configuration change.
   
   b. If somehow "9.a" above works, since both the controller-cards expect to be the master, none will power-down immediately as per Rule#1 above. Therefore the Chassis would need to be cold booted.

Now after the system reboots following "9.b", if we consider that the original problem of broken communication between the controller-cards persists, the configration of `config/power-admin-state=POWER_DISABLED` on one of the controller-cards wouldn't help because of Rule#5 above. However, it is expected that the implementation has other hardware/software means to gracefully handle the split-brain situation. The proposal here to allow for power-disable of a controller-card using configuration has no impact whatsoever on the implementations ability to handle split-brain like conditions. 
   
   ***Definition for Split-brain situation:*** 
   `In network routers/switches that use redundant controller-cards for high availability, a split-brain scenario occurs when the primary and secondary cards lose communication with each other, and both believe they should be the active controller. Probable causees for a Split-brain scenario can be configuration error or software bug preventing exchange of control messages between the cards breaking the communication`

   
   
