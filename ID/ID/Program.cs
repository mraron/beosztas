using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.IO;
using System.Diagnostics;
using PdfSharp.Pdf;
using PdfSharp.Drawing;
using PdfSharp.Drawing.Layout;

namespace RsaCryptoExample
{
    static class Program
    {
        static void Main()
        {
            string[] files = Directory.GetFiles("Nevsor", "*.txt");
            var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
            var random = new Random();
            foreach(string fn in files)
            {
                List<string> names = new List<string>();
                using (StreamReader sr = new StreamReader(fn, Encoding.GetEncoding("iso-8859-1")))
                {
                    while (sr.Peek() >= 0)
                    {
                        names.Add(sr.ReadLine());
                    }
                }
                int n = names.Count();
                string[] passwords = new string[n];
                for (int i = 0; i < n; i++)
                {
                    for (int j = 0; j < 10; j++)
                    {
                        passwords[i] += chars[random.Next(chars.Length)];
                    }
                    Console.WriteLine(passwords[i]);
                }
                PdfDocument document = new PdfDocument();
                string fnPdf = "";
                string str = "";
                for(int i = 0; i < fn.Length; i++)
                {
                    if(fn[i] == '.')
                    {
                        fnPdf += str;
                        str = "";
                    }
                    str += fn[i];
                }
                document.Info.Title = fnPdf;
                List<PdfPage> pages = new List<PdfPage>();
                List<XGraphics> gfx = new List<XGraphics>();
                List<XTextFormatter> tf = new List<XTextFormatter>();
                pages.Add(document.AddPage());
                gfx.Add(XGraphics.FromPdfPage(pages[0]));
                tf.Add(new XTextFormatter(gfx[0]));
                XFont titleFont = new XFont("Times New Roman", 32, XFontStyle.Bold);
                XFont font = new XFont("Times New Roman", 20, XFontStyle.Regular);
                tf[0].Alignment = XParagraphAlignment.Center;
                tf[0].DrawString(fnPdf.Split('\\')[fnPdf.Split('\\').Length - 1], titleFont, XBrushes.Black, new XRect(10, 10, pages[0].Width - 20, pages[0].Height - 40), XStringFormats.TopLeft);
                XPen pen = new XPen(XColors.Black, 1);
                gfx[0].DrawLine(pen, 0, 70, pages[0].Width, 70);
                int pageIndex = 0, newline = 80;
                for (int i = 0; i < n; i++)
                {
                    tf[pageIndex].Alignment = XParagraphAlignment.Left;
                    tf[pageIndex].DrawString(names[i], font, XBrushes.Black, new XRect(20, newline, pages[pageIndex].Width - 20, pages[pageIndex].Height - 40), XStringFormats.TopLeft);
                    tf[pageIndex].Alignment = XParagraphAlignment.Left;
                    tf[pageIndex].DrawString(passwords[i], font, XBrushes.Black, new XRect(pages[pageIndex].Width / 2 + 120, newline, pages[pageIndex].Width - 20, pages[pageIndex].Height - 40), XStringFormats.TopLeft);
                    gfx[pageIndex].DrawLine(pen, 0, newline - 10, pages[pageIndex].Width, newline - 10);
                    if (newline + 40 > pages[pageIndex].Height)
                    {
                        gfx[pageIndex].DrawLine(pen, 0, newline + 30, pages[pageIndex].Width, newline + 30);
                        pageIndex++;
                        pages.Add(document.AddPage());
                        gfx.Add(XGraphics.FromPdfPage(pages[pageIndex]));
                        tf.Add(new XTextFormatter(gfx[pageIndex]));
                        tf[pageIndex].Alignment = XParagraphAlignment.Center;
                        tf[pageIndex].DrawString(fnPdf.Split('\\')[fnPdf.Split('\\').Length - 1], titleFont, XBrushes.Black, new XRect(10, 10, pages[0].Width - 20, pages[0].Height - 10), XStringFormats.TopLeft);
                        gfx[pageIndex].DrawLine(pen, 0, 70, pages[pageIndex].Width, 70);
                        gfx[pageIndex - 1].DrawLine(pen, pages[pageIndex - 1].Width / 2 + 100, 70, pages[pageIndex - 1].Width / 2 + 100, newline - 10);
                        newline = 80;
                        tf[pageIndex].Alignment = XParagraphAlignment.Left;
                        tf[pageIndex].DrawString(names[i], font, XBrushes.Black, new XRect(20, newline, pages[pageIndex].Width - 20, pages[pageIndex].Height - 40), XStringFormats.TopLeft);
                        tf[pageIndex].Alignment = XParagraphAlignment.Left;
                        tf[pageIndex].DrawString(passwords[i], font, XBrushes.Black, new XRect(pages[pageIndex].Width / 2 + 120, newline, pages[pageIndex].Width - 20, pages[pageIndex].Height - 40), XStringFormats.TopLeft);
                        gfx[pageIndex].DrawLine(pen, 0, newline - 10, pages[pageIndex].Width, newline - 10);
                    }
                    newline += 40;
                }
                gfx[pageIndex].DrawLine(pen, 0, newline - 10, pages[pageIndex].Width, newline - 10);
                gfx[pageIndex].DrawLine(pen, pages[pageIndex].Width / 2 + 100, 70, pages[pageIndex].Width / 2 + 100, newline - 10);
                string filename = fnPdf + ".pdf";
                document.Save(filename);
                Process.Start(filename);
            }
        }
    }
}