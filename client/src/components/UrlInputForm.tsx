import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

interface UrlInputFormProps {
  onSubmit: (url: string) => void;
  autoFocus?: boolean;
  className?: string;
}

const formSchema = z.object({
  url: z.string().url({ message: "Please enter a valid URL." }),
});

export function UrlInputForm({ onSubmit, autoFocus, className }: UrlInputFormProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      url: "",
    },
  });

  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (autoFocus && inputRef.current) {
      inputRef.current.focus();
    }
  }, [autoFocus]);

  function handleSubmit(values: z.infer<typeof formSchema>) {
    onSubmit(values.url);
    form.reset();
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(handleSubmit)}
        className={`flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-2 w-full ${className}`}
      >
        <FormField
          control={form.control}
          name="url"
          render={({ field }) => (
            <FormItem className="flex-grow">
              <FormControl>
                <Input placeholder="https://example.com" {...field} ref={inputRef} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Analyze</Button>
      </form>
    </Form>
  );
}
