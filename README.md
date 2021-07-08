<p align="center">
  <img src="img/logo_small.png" alt="logo" width="250" align="middle"/>
</p>

*Nowadays space exploration is limited by technology and physics. Imagine finding yourself in a future where space travel is much easier. The goal of this project is to exploit serverless computing in order to find potentially habitable planets. A telescope will be in charge of spotting planets with characteristics similar to those of the earth, meanwhile a probe will be sent on that planet to obtain more accurate informations.*

-----

# Summary

[> Introduction](#Introduction)\
[> Architecture](#Architecture)\
[> Project structure](#Project-structure)\
[> Getting started](#Getting-started)

# Introduction

This project has been realized for the Serverless Computing for IoT class at "Università degli Studi di Salerno".

ProbeX takes place in a future in which interstellar travels are possible. The goal of the project is to discover potentially habitable planets exploiting the Serverless computing.

In the following sections will be explained how to achive the goal, which are the main technologies involved in this project and will be discussed all the solution design decisions.

# Architecture

<img src="img/architecture.png" alt="architecture"/>

ProbeX uses MQTT and AMPQ to keep in comunication all involved devices.

The project structure is composed by two main devices: the telescope, and the probe dock, each with a custom MQTT topic on which send messages.
When the Telescope finds a new planet, it sends a message on the MQTT topic *iot/telescope*. The handler *telescope-receiver* catch all telescope messages and finds out if the planet respects the habitability characteristics of the planet using a classifier (the current implementation of the classifier is just a mock).

Whenever the classifier returns as answer YES, the *telescope-receiver* sends a message on the topic *iot/probe*.

The probe dock catches that message and sends a probe on that planet. Each probe will communicate with ProbeX with the MQTT topic *iot/probe/receiver*.

An handler called *probe-receiver* will gain those informations on that topic, updating the *Monitor* status with all new informations catched by the probe.

Eather *telescope-receiver* or *probe-receiver* are implemented on *Nuclio*.

### QoS and Protocols

In this section will be explained all choises regarding the protocol choosed for each device, and the corrispective QoS (Quality of Service).

* **Telescope**: Since a Telescope is a powerful IoT device, the QoS is 2. Thats why it needs a reliable communication with the probe dock, otherwise more than one probe will be sent on the same planet.
* **Probe**: A probe will start to send information when it will arrive on the planet. The probe sends continous information, so we just need to guarantee that at least one message will arrive on our planet. That's why we just need QoS 1.
* **Monitor and Logger**: The monitor and the logger will communicate directly with the AMPQ protocol. Thats why a dedicated computer will be used as monitor and logger.

# Project Structure

- **devices**: contains the implementation of the simultated devices;
  - *probe.go*: the simulation device for the Probe Dock;
  - *telescope.go*: the simulation device for the Telescope;
- **probe**: defines the behaviour of **probe-receiver**;
  - *constants.go*: defines constants for the **probe-receiver**;
  - *consumer.go*: in this file is specified the function that will be executed on nuclio;
  - *logger.go*: contains all the utilities for logging on AMQP;
  - *persistence.go*: this file defines the communication of the probe-receiver with the monitor;
  - *function.yaml*: contains all the deploying informations of the function on nuclio;
- **telescope**: defines the behaviour of **telescope-receiver**;
  - *constants.go*: defines constants for the **probe-receiver**;
  - *consumer.go*: in this file is specified the function that will be executed on nuclio;
  - *logger.go*: contains all the utilities for logging on AMQP;
  - *classificator.go*: the mock classificator is defined in this file;
  - *function.yaml*: contains all the deploying informations of the function on nuclio;
- **monitor**: contains the implementation of the AMPQ monitor;
- **logging**: contains the implementation of the AMPQ logger;
- **.env**: contains the environment variables;

# Getting Started
