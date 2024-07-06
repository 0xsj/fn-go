import { useState } from "react";
import {
  Form,
  FormField,
  FormLabel,
  FormControl,
  FormMessage,
  FormItem,
} from "@/components/ui/form";
import React, { HtmlHTMLAttributes } from "react";
import { PasswordInput } from "@/components/custom/password-input";
import { Input } from "@/components/ui/input";
import { Link } from "react-router-dom";
import { Button } from "@/components/custom/button";
import { IconBrandFacebook, IconBrandGithub } from "@tabler/icons-react";
import { cn } from "@/utils/utils";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

interface Props extends HtmlHTMLAttributes<HTMLDivElement> {}

const formSchema = z.object({
  email: z.string().email(),
  password: z.string(),
});

export const SignInForm: React.FC<Props> = ({ className, ...props }) => {
  const [isLoading, setIsLoading] = useState<boolean>();
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onHandleSubmit = () => {
    return false;
  };

  return (
    <div className={cn("grid gap-6", className)} {...props}>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onHandleSubmit)}>
          <div className='grid gap-2'>
            <FormField
              control={form.control}
              name='email'
              render={({ field }) => (
                <FormItem className='space-y-1'>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input placeholder='name@example.com' {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='password'
              render={({ field }) => (
                <FormItem className='space-y-1'>
                  <div className='flex items-center justify-between'>
                    <FormLabel>Password</FormLabel>
                    <Link
                      to='/forgot-password'
                      className='text-sm font-medium text-muted-foreground hover:opacity-75'
                    >
                      Forgot password?
                    </Link>
                  </div>
                  <FormControl>
                    <PasswordInput placeholder='********' {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button className='mt-2' loading={false}>
              Login
            </Button>

            <div className='relative my-2'>
              <div className='absolute inset-0 flex items-center'>
                <span className='w-full border-t' />
              </div>
              <div className='relative flex justify-center text-xs uppercase'>
                <span className='bg-background px-2 text-muted-foreground'>
                  Or continue with
                </span>
              </div>
            </div>

            <div className='flex items-center gap-2'>
              <Button
                variant='outline'
                className='w-full'
                type='button'
                loading={isLoading}
                leftSection={<IconBrandGithub className='h-4 w-4' />}
              >
                GitHub
              </Button>
              <Button
                variant='outline'
                className='w-full'
                type='button'
                loading={isLoading}
                leftSection={<IconBrandFacebook className='h-4 w-4' />}
              >
                Facebook
              </Button>
            </div>
          </div>
        </form>
      </Form>
    </div>
  );
};
