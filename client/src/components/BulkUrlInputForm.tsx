import { zodResolver } from "@hookform/resolvers/zod";
import { ClipboardList } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";

interface BulkUrlInputFormProps {
  onSubmit: (urls: string[]) => void;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  displayMode?: "full" | "icon";
}

const formSchema = z.object({
  urls: z
    .string()
    .min(1, { message: "Please enter at least one URL." })
    .refine(
      (value) => {
        const urls = value.split("\n").filter(Boolean);
        return urls.every((url) => z.string().url().safeParse(url).success);
      },
      { message: "Please enter a valid URL per line." },
    ),
});

export function BulkUrlInputForm({
  onSubmit,
  open,
  onOpenChange,
  displayMode = "full",
}: BulkUrlInputFormProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { urls: "" },
  });

  function handleSubmit(values: z.infer<typeof formSchema>) {
    const urls = values.urls.split("\n").filter(Boolean);
    onSubmit(urls);
    form.reset();
    onOpenChange(false);
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        {displayMode === "icon" ? (
          <Button variant="outline" size="icon">
            <ClipboardList className="h-4 w-4" />
          </Button>
        ) : (
          <Button variant="outline" className="w-full flex items-center gap-2">
            <ClipboardList className="h-4 w-4" />
            Bulk Analyze
          </Button>
        )}
      </DialogTrigger>

      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Bulk URL Analysis</DialogTitle>
          <DialogDescription>
            Enter one URL per line to analyze multiple websites at once.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="urls"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Textarea
                      placeholder={`https://example.com\nhttps://google.com`}
                      {...field}
                      rows={5}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="submit">Analyze</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
