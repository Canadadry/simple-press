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
} from "@radix-ui/themes";
import type { DynamicFormUI } from "./render";
import {
  TrashIcon,
  ChevronUpIcon,
  ChevronDownIcon,
  Pencil1Icon,
  TextAlignLeftIcon,
} from "@radix-ui/react-icons";

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

    FormObject: ({ label, children }) => {
      return (
        <Card mb="2">
          {label}
          <Box mt="2">{children}</Box>
        </Card>
      );
    },

    FormInput: ({ label, name, inputType, value, setData }) => (
      <TextField.Root
        mb="4"
        data-testid={`input-${name}`}
        defaultValue={value}
        type={inputType as "text" | "number"}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
          setData(name, e.target.value)
        }
      >
        <TextField.Slot>{label}</TextField.Slot>
      </TextField.Root>
    ),

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
