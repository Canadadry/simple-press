import React, { useState } from "react";
import {
  Card,
  Flex,
  Box,
  Text,
  TextField,
  Checkbox,
  Button,
  Spinner,
  TextArea,
} from "@radix-ui/themes";
import type { DynamicFormUI } from "./render";
import {
  TrashIcon,
  ChevronUpIcon,
  ChevronDownIcon,
  Pencil1Icon,
  TextAlignLeftIcon,
} from "@radix-ui/react-icons";
import * as Accordion from "@radix-ui/react-accordion";

export function makeRadixUI(maxWidth: number): DynamicFormUI {
  return {
    Form: ({
      label,
      children,
      mode,
      setMode,
      saving,
      onSave,
      onDelete,
      onUp,
      onDown,
    }) => (
      <Card>
        <Flex gap="2" justify="between">
          <Text as="div" size="2" mb="2" weight="bold" color="indigo">
            {label}
            {mode === "form" ? (
              <Pencil1Icon onClick={() => setMode("json")}></Pencil1Icon>
            ) : (
              <TextAlignLeftIcon
                onClick={() => setMode("form")}
              ></TextAlignLeftIcon>
            )}
          </Text>
          <Flex gap="2">
            <ChevronUpIcon onClick={onUp}></ChevronUpIcon>
            <ChevronDownIcon onClick={onDown}></ChevronDownIcon>
          </Flex>
        </Flex>
        <Box
          data-testid="form"
          style={{
            maxWidth: maxWidth,
          }}
        >
          {children}
          <Flex gap="2" justify="between">
            <Button
              tabIndex={1}
              size="2"
              variant="outline"
              color="crimson"
              onClick={onDelete}
            >
              <TrashIcon />
              Delete
            </Button>
            <Button
              tabIndex={1}
              size="2"
              onClick={onSave}
              disabled={saving != "touched"}
            >
              {saving == "saving" ? <Spinner /> : "Save"}
            </Button>
          </Flex>
        </Box>
      </Card>
    ),

    FormObject: ({ label, children }) => (
      <Card mb="2">
        <Accordion.Root type="single" collapsible>
          <Accordion.Item value={label}>
            <Accordion.Header style={{ margin: 0, padding: 0 }}>
              <Accordion.Trigger
                style={{
                  all: "unset",
                  width: "100%",
                  cursor: "pointer",
                }}
              >
                <Flex justify="between">
                  <Text size="2" weight="bold">
                    {label}
                  </Text>
                  <ChevronDownIcon />
                </Flex>
              </Accordion.Trigger>
            </Accordion.Header>

            <Accordion.Content>
              <Box mt="2" data-testid={`object-${label}`}>
                {children}
              </Box>
            </Accordion.Content>
          </Accordion.Item>
        </Accordion.Root>
      </Card>
    ),

    FormInput: ({ label, name, inputType, value, setData }) => {
      const part = label.split("__");
      if (part[1] && part[1].startsWith("ta")) {
        let row = 1;
        try {
          row = parseInt(part[1].slice(2));
        } catch {
          // nothing
        }
        return (
          <Box position="relative">
            <TextArea
              mb="4"
              style={{ paddingTop: "2rem" }}
              spellCheck={false}
              variant="surface"
              rows={row}
              value={value}
              onChange={(e) => {
                setData(name, e.target.value);
              }}
            />
            <Box position="absolute" m="2" top="0" left="0" right="0">
              <strong>{part[0]}</strong>
            </Box>
          </Box>
        );
      }
      return (
        <TextField.Root
          mb="4"
          data-testid={`input-${name}`}
          defaultValue={value}
          type={inputType as "text" | "number"}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setData(name, e.target.value)
          }
        >
          <TextField.Slot>{part[0]}</TextField.Slot>
        </TextField.Root>
      );
    },

    FormCheckBox: ({ label, name, checked, setData }) => {
      const [localChecked, setLocalChecked] = useState(checked);
      return (
        <Text as="label" size="2" data-testid={`checkbox-${name}`}>
          <Flex gap="2">
            <Checkbox
              mb="2"
              checked={localChecked}
              onCheckedChange={(c) => {
                const value = c === "indeterminate" ? false : c;
                setLocalChecked(value);
                setData(name, String(value));
              }}
            />
            {label}
          </Flex>
        </Text>
      );
    },
  };
}
